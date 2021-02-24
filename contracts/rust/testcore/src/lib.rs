// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

use consts::*;
use testcore::*;
use wasmlib::*;

mod testcore;
mod consts;

#[no_mangle]
fn on_load() {
    let exports = ScExports::new();
    exports.add_func(FUNC_CALL_ON_CHAIN, func_call_on_chain);
    exports.add_func(FUNC_CHECK_CONTEXT_FROM_FULL_EP, func_check_context_from_full_ep);
    exports.add_func(FUNC_DO_NOTHING, func_do_nothing);
    exports.add_func(FUNC_INIT, func_init);
    exports.add_func(FUNC_PASS_TYPES_FULL, func_pass_types_full);
    exports.add_func(FUNC_RUN_RECURSION, func_run_recursion);
    exports.add_func(FUNC_SEND_TO_ADDRESS, func_send_to_address);
    exports.add_func(FUNC_SET_INT, func_set_int);
    exports.add_func(FUNC_GET_MINTED_SUPPLY, func_get_minted_supply);
    exports.add_func(FUNC_TEST_CALL_PANIC_FULL_EP, func_test_call_panic_full_ep);
    exports.add_func(FUNC_TEST_CALL_PANIC_VIEW_EPFROM_FULL, func_test_call_panic_view_epfrom_full);
    exports.add_func(FUNC_TEST_CHAIN_OWNER_IDFULL, func_test_chain_owner_idfull);
    exports.add_func(FUNC_TEST_CONTRACT_IDFULL, func_test_contract_idfull);
    exports.add_func(FUNC_TEST_EVENT_LOG_DEPLOY, func_test_event_log_deploy);
    exports.add_func(FUNC_TEST_EVENT_LOG_EVENT_DATA, func_test_event_log_event_data);
    exports.add_func(FUNC_TEST_EVENT_LOG_GENERIC_DATA, func_test_event_log_generic_data);
    exports.add_func(FUNC_TEST_PANIC_FULL_EP, func_test_panic_full_ep);
    exports.add_func(FUNC_WITHDRAW_TO_CHAIN, func_withdraw_to_chain);
    exports.add_view(VIEW_CHECK_CONTEXT_FROM_VIEW_EP, view_check_context_from_view_ep);
    exports.add_view(VIEW_FIBONACCI, view_fibonacci);
    exports.add_view(VIEW_GET_COUNTER, view_get_counter);
    exports.add_view(VIEW_GET_INT, view_get_int);
    exports.add_func(FUNC_INC_COUNTER, func_inc_counter);
    exports.add_view(VIEW_JUST_VIEW, view_just_view);
    exports.add_view(VIEW_PASS_TYPES_VIEW, view_pass_types_view);
    exports.add_view(VIEW_TEST_CALL_PANIC_VIEW_EPFROM_VIEW, view_test_call_panic_view_epfrom_view);
    exports.add_view(VIEW_TEST_CHAIN_OWNER_IDVIEW, view_test_chain_owner_idview);
    exports.add_view(VIEW_TEST_CONTRACT_IDVIEW, view_test_contract_idview);
    exports.add_view(VIEW_TEST_PANIC_VIEW_EP, view_test_panic_view_ep);
    exports.add_view(VIEW_TEST_SANDBOX_CALL, view_test_sandbox_call);
}
