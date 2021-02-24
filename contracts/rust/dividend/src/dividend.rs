// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

use wasmlib::*;

use crate::*;
use crate::types::*;

pub fn func_divide(ctx: &ScFuncContext) {
    ctx.log("dividend.divide");
    let amount = ctx.balances().balance(&ScColor::IOTA);
    if amount == 0 {
        ctx.panic("Nothing to divide");
    }
    let state = ctx.state();
    let total_factor = state.get_int(VAR_TOTAL_FACTOR);
    let total = total_factor.value();
    let members = state.get_bytes_array(VAR_MEMBERS);
    let mut parts = 0_i64;
    let size = members.length();
    for i in 0..size {
        let m = Member::from_bytes(&members.get_bytes(i).value());
        let part = amount * m.factor / total;
        if part != 0 {
            parts += part;
            ctx.transfer_to_address(&m.address, ScTransfers::new(&ScColor::IOTA, part));
        }
    }
    if parts != amount {
        // note we truncated the calculations down to the nearest integer
        // there could be some small remainder left in the contract, but
        // that will be picked up in the next round as part of the balance
        let remainder = amount - parts;
        ctx.log(&("Remainder in contract: ".to_string() + &remainder.to_string()));
    }
    ctx.log("dividend.divide ok");
}

pub fn func_member(ctx: &ScFuncContext) {
    ctx.log("dividend.member");
    // only creator can add members
    ctx.require(ctx.caller() == ctx.contract_creator(), "no permission");

    let p = ctx.params();
    let address = p.get_address(PARAM_ADDRESS);
    let factor = p.get_int(PARAM_FACTOR);
    ctx.require(address.exists(), "missing mandatory address");
    ctx.require(factor.exists(), "missing mandatory factor");
    let member = Member {
        address: address.value(),
        factor: factor.value(),
    };
    let state = ctx.state();
    let total_factor = state.get_int(VAR_TOTAL_FACTOR);
    let mut total = total_factor.value();
    let members = state.get_bytes_array(VAR_MEMBERS);
    let size = members.length();
    for i in 0..size {
        let m = Member::from_bytes(&members.get_bytes(i).value());
        if m.address == member.address {
            total -= m.factor;
            total += member.factor;
            total_factor.set_value(total);
            members.get_bytes(i).set_value(&member.to_bytes());
            ctx.log(&("Updated: ".to_string() + &member.address.to_string()));
            return;
        }
    }
    total += member.factor;
    total_factor.set_value(total);
    members.get_bytes(size).set_value(&member.to_bytes());
    ctx.log(&("Appended: ".to_string() + &member.address.to_string()));
    ctx.log("dividend.member ok");
}
