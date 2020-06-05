package cluster

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"text/template"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/transaction"
	nodeapi "github.com/iotaledger/goshimmer/packages/waspconn/apilib"
	"github.com/iotaledger/goshimmer/packages/waspconn/utxodb"
	waspapi "github.com/iotaledger/wasp/packages/apilib"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/sctransaction"
	"github.com/iotaledger/wasp/packages/sctransaction/origin"
	"github.com/iotaledger/wasp/packages/util"
)

type SmartContractKeys struct {
	Address  string
	DKShares []string // [node index]
}

type SmartContractConfig struct {
	Description  string `json:"description"`
	OwnerAddress string `json:"ownerAddress"`
	Nodes        []int  `json:"nodes"`
	Quorum       int    `json:"quorum"`
}

type ClusterConfig struct {
	Nodes []struct {
		BindAddress string `json:"bindAddress"`
		PeeringPort int    `json:"PeeringPort"`
	} `json:"nodes"`
	Goshimmer      string `json:"goshimmer"`
	SmartContracts []SmartContractConfig
}

type Cluster struct {
	Config     *ClusterConfig
	ConfigPath string // where the cluster configuration is stored - read only
	DataPath   string // where the cluster's volatile data lives
	Started    bool
	cmds       []*exec.Cmd
}

func readConfig(configPath string) (*ClusterConfig, error) {
	data, err := ioutil.ReadFile(path.Join(configPath, "cluster.json"))
	if err != nil {
		return nil, err
	}

	config := &ClusterConfig{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func New(configPath string, dataPath string) (*Cluster, error) {
	config, err := readConfig(configPath)
	if err != nil {
		return nil, err
	}
	return &Cluster{
		Config:     config,
		ConfigPath: configPath,
		DataPath:   dataPath,
		cmds:       make([]*exec.Cmd, 0),
	}, nil
}

func (cluster *Cluster) readKeysConfig() ([]SmartContractKeys, error) {
	data, err := ioutil.ReadFile(cluster.ConfigKeysPath())
	if err != nil {
		return nil, err
	}

	config := make([]SmartContractKeys, 0)
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (cluster *Cluster) JoinConfigPath(s string) string {
	return path.Join(cluster.ConfigPath, s)
}

func (cluster *Cluster) ConfigTemplatePath() string {
	return cluster.JoinConfigPath("wasp-config-template.json")
}

func (cluster *Cluster) ConfigKeysPath() string {
	return cluster.JoinConfigPath("keys.json")
}

func (cluster *Cluster) NodeDataPath(i int) string {
	return path.Join(cluster.DataPath, strconv.Itoa(i))
}

func (cluster *Cluster) JoinNodeDataPath(i int, s string) string {
	return path.Join(cluster.NodeDataPath(i), s)
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Init creates in DataPath a directory with config.json for each node
func (cluster *Cluster) Init(resetDataPath bool) error {
	exists, err := fileExists(cluster.DataPath)
	if err != nil {
		return err
	}
	if exists {
		if !resetDataPath {
			return errors.New(fmt.Sprintf("%s directory exists", cluster.DataPath))
		}
		err = os.RemoveAll(cluster.DataPath)
		if err != nil {
			return err
		}
	}

	configTmpl, err := template.ParseFiles(cluster.ConfigTemplatePath())
	if err != nil {
		return err
	}

	for i, nodeConfig := range cluster.Config.Nodes {
		nodePath := cluster.NodeDataPath(i)
		fmt.Printf("Initializing node configuration at %s\n", nodePath)

		err := os.MkdirAll(nodePath, os.ModePerm)
		if err != nil {
			return err
		}

		f, err := os.Create(cluster.JoinNodeDataPath(i, "config.json"))
		if err != nil {
			return err
		}
		defer f.Close()
		err = configTmpl.Execute(f, &nodeConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

func logNode(i int, scanner *bufio.Scanner, initString string, initOk chan bool) {
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if !found && strings.Contains(line, initString) {
			initOk <- true
			found = true
		}
		fmt.Printf("[wasp %d] %s\n", i, line)
	}
}

// Start launches all wasp nodes in the cluster, each running in its own directory
func (cluster *Cluster) Start() error {
	exists, err := fileExists(cluster.DataPath)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New(fmt.Sprintf("Data path %s does not exist", cluster.DataPath))
	}

	err = cluster.start()
	if err != nil {
		return err
	}

	err = cluster.importKeys()
	if err != nil {
		return err
	}
	cluster.Started = true
	return nil
}

func (cluster *Cluster) start() error {
	initOk := make(chan bool, len(cluster.Config.Nodes))

	for i, _ := range cluster.Config.Nodes {
		cmd := exec.Command("wasp")
		cmd.Dir = cluster.NodeDataPath(i)
		pipe, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(pipe)
		err = cmd.Start()
		if err != nil {
			return err
		}
		cluster.cmds = append(cluster.cmds, cmd)

		go logNode(i, scanner, "WebAPI started", initOk)
	}

	for i := 0; i < len(cluster.Config.Nodes); i++ {
		<-initOk
	}
	return nil
}

func (cluster *Cluster) importKeys() error {
	exists, err := fileExists(cluster.ConfigKeysPath())
	if err != nil {
		return err
	}
	if !exists {
		// nothing to do
		return nil
	}

	keys, err := cluster.readKeysConfig()
	if err != nil {
		return err
	}

	for _, scKeys := range keys {
		fmt.Printf("Importing DKShares for account %s...\n", scKeys.Address)
		for nodeIndex, dks := range scKeys.DKShares {
			err := waspapi.ImportDKShare(cluster.Config.Nodes[nodeIndex].BindAddress, dks)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Stop sends an interrupt signal to all nodes and waits for them to exit
func (cluster *Cluster) Stop() {
	for _, node := range cluster.Config.Nodes {
		fmt.Printf("Sending shutdown to %s\n", node.BindAddress)
		err := waspapi.Shutdown(node.BindAddress)
		if err != nil {
			fmt.Println(err)
		}
	}
	cluster.Wait()
}

// Wait blocks until all nodes exit
func (cluster *Cluster) Wait() {
	for _, cmd := range cluster.cmds {
		err := cmd.Wait()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (cluster *Cluster) Hosts() []string {
	hosts := make([]string, 0)
	for _, node := range cluster.Config.Nodes {
		hosts = append(hosts, node.BindAddress)
	}
	return hosts
}

func (cluster *Cluster) Committee(sc *SmartContractConfig) ([]string, error) {
	committee := make([]string, 0)
	for _, i := range sc.Nodes {
		if i < 0 || i > len(cluster.Config.Nodes)-1 {
			return nil, errors.New(fmt.Sprintf("Node index out of bounds in smart contract committee configuration: %d", i))
		}
		committee = append(committee, cluster.Config.Nodes[i].BindAddress)
	}
	return committee, nil

}

func (cluster *Cluster) GenerateDKSets() error {
	keysFile := cluster.ConfigKeysPath()
	exists, err := fileExists(keysFile)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("dk sets already generated in keys.json")
	}

	keys := make([]SmartContractKeys, 0)

	for _, sc := range cluster.Config.SmartContracts {
		committee, err := cluster.Committee(&sc)
		if err != nil {
			return err
		}
		addr, err := waspapi.GenerateNewDistributedKeySet(
			committee,
			uint16(len(committee)),
			uint16(sc.Quorum),
		)
		if err != nil {
			return err
		}

		fmt.Printf("Generated key set for SC with address %s\n", addr)

		dkShares := make([]string, 0)
		for _, host := range cluster.Hosts() {
			dks, err := waspapi.ExportDKShare(host, addr)
			if err != nil {
				return err
			}
			dkShares = append(dkShares, dks)
		}

		keys = append(keys, SmartContractKeys{
			Address:  addr.String(),
			DKShares: dkShares,
		})
	}
	buf, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(keysFile, buf, 0644)
}

func (cluster *Cluster) CreateOriginTx() error {
	keys, err := cluster.readKeysConfig()
	if err != nil {
		return err
	}

	for scIndex, sc := range cluster.Config.SmartContracts {
		tx, err := cluster.createOriginTx(&sc, &keys[scIndex])
		if err != nil {
			return err
		}

		fmt.Printf("Posting origin tx for SC %d: txid %s\n", scIndex, tx.ID().String())

		err = nodeapi.PostTransaction(cluster.Config.Goshimmer, tx.Transaction)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cluster *Cluster) createOriginTx(sc *SmartContractConfig, keys *SmartContractKeys) (*sctransaction.Transaction, error) {
	ownerAddress, err := address.FromBase58(sc.OwnerAddress)
	if err != nil {
		return nil, err
	}
	scAddr, err := address.FromBase58(keys.Address)
	if err != nil {
		return nil, err
	}

	inputTxId, inputBalances, err := cluster.selectInputFromAvailableOutputs(&ownerAddress)
	if err != nil {
		return nil, err
	}

	originTx, err := origin.NewOriginTransaction(origin.NewOriginTransactionParams{
		NewOriginParams: origin.NewOriginParams{
			Address:      scAddr,
			OwnerAddress: ownerAddress,
			ProgramHash:  *hashing.HashStrings(sc.Description), // TODO
		},
		Input:          *inputTxId,
		InputBalances:  inputBalances,
		InputColor:     balance.ColorIOTA,
		OwnerSigScheme: utxodb.GetSigScheme(ownerAddress),
	})
	if err != nil {
		return nil, err
	}
	return originTx, nil
}

func (cluster *Cluster) selectInputFromAvailableOutputs(ownerAddress *address.Address) (input *transaction.OutputID, inputBalances []*balance.Balance, err error) {
	allOutputs, err := nodeapi.GetAccountOutputs(cluster.Config.Goshimmer, ownerAddress)
	if err != nil {
		return
	}

	outputs := util.SelectMinimumOutputs(allOutputs, balance.ColorIOTA, 1)
	if len(outputs) == 0 {
		err = fmt.Errorf("Not enough outputs for 1 iota!")
		return
	}

	// len(outputs) should be 1
	for oid, v := range outputs {
		input = &oid
		inputBalances = v
		break
	}
	return
}