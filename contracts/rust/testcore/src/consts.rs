// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

#![allow(dead_code)]

use wasmlib::*;

pub const SC_NAME: &str = "testcore";
pub const SC_DESCRIPTION: &str = "Core test for ISCP wasmlib Rust/Wasm library";
pub const SC_HNAME: ScHname = ScHname(0x370d33ad);

pub const PARAM_ADDRESS: &str = "address";
pub const PARAM_AGENT_ID: &str = "agentID";
pub const PARAM_CALLER: &str = "caller";
pub const PARAM_CHAIN_ID: &str = "chainid";
pub const PARAM_CHAIN_OWNER_ID: &str = "chainOwnerID";
pub const PARAM_CONTRACT_CREATOR: &str = "contractCreator";
pub const PARAM_CONTRACT_ID: &str = "contractID";
pub const PARAM_COUNTER: &str = "counter";
pub const PARAM_MINTED_SUPPLY: &str = "mintedSupply";
pub const PARAM_HASH: &str = "Hash";
pub const PARAM_HNAME: &str = "Hname";
pub const PARAM_HNAME_CONTRACT: &str = "hnameContract";
pub const PARAM_HNAME_EP: &str = "hnameEP";
pub const PARAM_HNAME_ZERO: &str = "Hname-0";
pub const PARAM_INT64: &str = "int64";
pub const PARAM_INT64_ZERO: &str = "int64-0";
pub const PARAM_INT_VALUE: &str = "intParamValue";
pub const PARAM_NAME: &str = "intParamName";
pub const PARAM_STRING: &str = "string";
pub const PARAM_STRING_ZERO: &str = "string-0";

pub const VAR_COUNTER: &str = "counter";
pub const VAR_HNAME_EP: &str = "hnameEP";

pub const FUNC_CALL_ON_CHAIN: &str = "callOnChain";
pub const FUNC_CHECK_CONTEXT_FROM_FULL_EP: &str = "checkContextFromFullEP";
pub const FUNC_DO_NOTHING: &str = "doNothing";
pub const FUNC_INIT: &str = "init";
pub const FUNC_PASS_TYPES_FULL: &str = "passTypesFull";
pub const FUNC_RUN_RECURSION: &str = "runRecursion";
pub const FUNC_SEND_TO_ADDRESS: &str = "sendToAddress";
pub const FUNC_SET_INT: &str = "setInt";
pub const FUNC_GET_MINTED_SUPPLY: &str = "getMintedSupply";
pub const FUNC_TEST_CALL_PANIC_FULL_EP: &str = "testCallPanicFullEP";
pub const FUNC_TEST_CALL_PANIC_VIEW_EPFROM_FULL: &str = "testCallPanicViewEPFromFull";
pub const FUNC_TEST_CHAIN_OWNER_IDFULL: &str = "testChainOwnerIDFull";
pub const FUNC_TEST_CONTRACT_IDFULL: &str = "testContractIDFull";
pub const FUNC_TEST_EVENT_LOG_DEPLOY: &str = "testEventLogDeploy";
pub const FUNC_TEST_EVENT_LOG_EVENT_DATA: &str = "testEventLogEventData";
pub const FUNC_TEST_EVENT_LOG_GENERIC_DATA: &str = "testEventLogGenericData";
pub const FUNC_TEST_PANIC_FULL_EP: &str = "testPanicFullEP";
pub const FUNC_WITHDRAW_TO_CHAIN: &str = "withdrawToChain";
pub const VIEW_CHECK_CONTEXT_FROM_VIEW_EP: &str = "checkContextFromViewEP";
pub const VIEW_FIBONACCI: &str = "fibonacci";
pub const FUNC_INC_COUNTER: &str = "incCounter";
pub const VIEW_GET_COUNTER: &str = "getCounter";
pub const VIEW_GET_INT: &str = "getInt";
pub const VIEW_JUST_VIEW: &str = "justView";
pub const VIEW_PASS_TYPES_VIEW: &str = "passTypesView";
pub const VIEW_TEST_CALL_PANIC_VIEW_EPFROM_VIEW: &str = "testCallPanicViewEPFromView";
pub const VIEW_TEST_CHAIN_OWNER_IDVIEW: &str = "testChainOwnerIDView";
pub const VIEW_TEST_CONTRACT_IDVIEW: &str = "testContractIDView";
pub const VIEW_TEST_PANIC_VIEW_EP: &str = "testPanicViewEP";
pub const VIEW_TEST_SANDBOX_CALL: &str = "testSandboxCall";

pub const HFUNC_CALL_ON_CHAIN: ScHname = ScHname(0x95a3d123);
pub const HFUNC_CHECK_CONTEXT_FROM_FULL_EP: ScHname = ScHname(0xa56c24ba);
pub const HFUNC_DO_NOTHING: ScHname = ScHname(0xdda4a6de);
pub const HFUNC_INIT: ScHname = ScHname(0x1f44d644);
pub const HFUNC_PASS_TYPES_FULL: ScHname = ScHname(0x733ea0ea);
pub const HFUNC_RUN_RECURSION: ScHname = ScHname(0x833425fd);
pub const HFUNC_SEND_TO_ADDRESS: ScHname = ScHname(0x63ce4634);
pub const HFUNC_SET_INT: ScHname = ScHname(0x62056f74);
pub const HFUNC_TEST_CALL_PANIC_FULL_EP: ScHname = ScHname(0x4c878834);
pub const HFUNC_TEST_CALL_PANIC_VIEW_EPFROM_FULL: ScHname = ScHname(0xfd7e8c1d);
pub const HFUNC_TEST_CHAIN_OWNER_IDFULL: ScHname = ScHname(0x2aff1167);
pub const HFUNC_TEST_CONTRACT_IDFULL: ScHname = ScHname(0x95934282);
pub const HFUNC_TEST_EVENT_LOG_DEPLOY: ScHname = ScHname(0x96ff760a);
pub const HFUNC_TEST_EVENT_LOG_EVENT_DATA: ScHname = ScHname(0x0efcf939);
pub const HFUNC_TEST_EVENT_LOG_GENERIC_DATA: ScHname = ScHname(0x6a16629d);
pub const HFUNC_TEST_PANIC_FULL_EP: ScHname = ScHname(0x24fdef07);
pub const HFUNC_WITHDRAW_TO_CHAIN: ScHname = ScHname(0x437bc026);
pub const HVIEW_CHECK_CONTEXT_FROM_VIEW_EP: ScHname = ScHname(0x88ff0167);
pub const HVIEW_FIBONACCI: ScHname = ScHname(0x7940873c);
pub const HVIEW_GET_COUNTER: ScHname = ScHname(0xb423e607);
pub const HVIEW_GET_INT: ScHname = ScHname(0x1887e5ef);
pub const HVIEW_JUST_VIEW: ScHname = ScHname(0x33b8972e);
pub const HVIEW_PASS_TYPES_VIEW: ScHname = ScHname(0x1a5b87ea);
pub const HVIEW_TEST_CALL_PANIC_VIEW_EPFROM_VIEW: ScHname = ScHname(0x91b10c99);
pub const HVIEW_TEST_CHAIN_OWNER_IDVIEW: ScHname = ScHname(0x26586c33);
pub const HVIEW_TEST_CONTRACT_IDVIEW: ScHname = ScHname(0x28a02913);
pub const HVIEW_TEST_PANIC_VIEW_EP: ScHname = ScHname(0x22bc4d72);
pub const HVIEW_TEST_SANDBOX_CALL: ScHname = ScHname(0x42d72b63);
