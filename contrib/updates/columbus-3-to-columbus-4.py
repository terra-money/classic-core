#!/usr/bin/env python3

import argparse
import json
import sys


def init_default_argument_parser(prog_desc, default_chain_id,
                                 default_genesis_time, defalut_halt_height, default_pretty):
    parser = argparse.ArgumentParser(description=prog_desc)
    parser.add_argument(
        'exported_genesis',
        help='exported genesis.json file',
        type=argparse.FileType('r'), default=sys.stdin,
    )
    parser.add_argument('--chain-id', type=str, default=default_chain_id)
    parser.add_argument('--genesis-time', type=str,
                        default=default_genesis_time)
    parser.add_argument('--halt-height', type=str,
                        default=defalut_halt_height)
    parser.add_argument('--pretty', type=bool,
                        default=default_pretty)
    return parser


def main(argument_parser, process_genesis_func):
    args = argument_parser.parse_args()
    if args.chain_id.strip() == '':
        sys.exit('chain-id required')

    genesis = json.loads(args.exported_genesis.read())
    genesis = process_genesis_func(genesis=genesis, parsed_args=args,)

    if args.pretty:
        raw_genesis = json.dumps(genesis, indent=4, sort_keys=True)
    else:
        raw_genesis = json.dumps(genesis, indent=None,
                                 sort_keys=False, separators=(',', ':'))

    print(raw_genesis)


def process_raw_genesis(genesis, parsed_args):

    # Auth Module Migration
    accs = []
    for acc in genesis['app_state']['accounts']:
        accType = 'core/Account'
        if acc['vesting_schedules']:
            # Vesting Acc
            accType = 'core/LazyGradedVestingAccount'
        else:
            if acc['module_name']:
                # Module Acc
                accType = 'supply/ModuleAccount'
                acc['name'] = acc['module_name']
                acc['permissions'] = acc['module_permissions']

            # Normal Acc
            del acc['delegated_free']
            del acc['delegated_vesting']
            del acc['end_time']
            del acc['original_vesting']
            del acc['vesting_schedules']

        del acc['start_time']
        del acc['module_name']
        del acc['module_permissions']

        accs.append({
            'type': accType,
            'value': acc,
        })

    del genesis['app_state']['accounts']
    genesis['app_state']['auth'] = {
        'accounts': accs,
        'params': genesis['app_state']['auth']['params']
    }

    # Distribution Module Migration
    genesis['app_state']['distribution']['params'] = {
        'community_tax': genesis['app_state']['distribution']['community_tax'],
        'base_proposer_reward': genesis['app_state']['distribution']['base_proposer_reward'],
        'bonus_proposer_reward': genesis['app_state']['distribution']['bonus_proposer_reward'],
        'withdraw_addr_enabled': genesis['app_state']['distribution']['withdraw_addr_enabled'],
    }

    # Staking Module Migration
    vals = []
    for val in genesis['app_state']['staking']['validators']:
        val['description']['security_contact'] = ''
        vals.append(val)

    genesis['app_state']['staking']['validators'] = vals

    # Slahing Module Migration to Evidence Module
    genesis['app_state']['evidence'] = {
        'params': {
            'max_evidence_age': genesis['app_state']['slashing']['params']['max_evidence_age'],
        },
        'evidence': [],
    }

    del genesis['app_state']['slashing']['params']['max_evidence_age']

    # Market Module Migration
    tobin_taxes = {}
    for item in genesis['app_state']['market']['params']['illiquid_tobin_tax_list']:
        tobin_taxes[item['denom']] = item['tax_rate']

    default_tobin_tax = genesis['app_state']['market']['params']['tobin_tax']
    del genesis['app_state']['market']['params']['tobin_tax']
    del genesis['app_state']['market']['params']['illiquid_tobin_tax_list']

    # Oracle Module Migration
    genesis['app_state']['oracle']['aggregate_exchange_rate_prevotes'] = []
    genesis['app_state']['oracle']['aggregate_exchange_rate_votes'] = []
    genesis['app_state']['oracle']['tobin_taxes'] = {}

    whitelist = []
    for item in genesis['app_state']['oracle']['params']['whitelist']:
        tobin_tax = default_tobin_tax
        if item in tobin_taxes:
            tobin_tax = tobin_taxes[item]

        whitelist.append({
            'name': item,
            'tobin_tax': tobin_tax,
        })

    genesis['app_state']['oracle']['params']['whitelist'] = whitelist

    # Treasury Module Migration
    genesis['app_state']['treasury']['cumulated_height'] = parsed_args.halt_height

    # Migrate Tendermint Consensus Param
    genesis['consensus_params']['evidence'] = {
        'max_age_num_blocks': genesis['consensus_params']['evidence']['max_age'],
        'max_age_duration': '172800000000000',
    }

    # Mint module
    genesis['app_state']['mint'] = {
        'minter': {
            'inflation': '0.0',
            'annual_provisions': '0.0',
        },
        'params': {
            'mint_denom': 'uluna',
            'inflation_rate_change': '0.0',
            'inflation_max': '0.0',
            'inflation_min': '0.0',
            'goal_bonded': '0.67',
            # uint64(60 * 60 * 8766 / 5), // assuming 5 second block times
            'blocks_per_year': '6311520',
        }
    }

    # MsgAuth module
    genesis['app_state']['msgauth'] = {
        'authorization_entries': []
    }

    # Wasm module
    genesis['app_state']['wasm'] = {
        'codes': [],
        'contracts': [],
        'last_code_id': '0',
        'last_instance_id': '0',
        'params': {
            "max_contract_gas": "100000000",
            "max_contract_msg_size": "10240",   # 10KB
            "max_contract_size": "512000"       # 500KB
        }
    }

    # Set new chain ID and genesis start time
    genesis['chain_id'] = parsed_args.chain_id.strip()
    genesis['genesis_time'] = parsed_args.genesis_time

    return genesis


if __name__ == '__main__':
    parser = init_default_argument_parser(
        prog_desc='Convert genesis.json for columbus-4',
        default_chain_id='columbus-4',
        default_genesis_time='2020-10-03T18:00:00Z',
        defalut_halt_height='3820000',
        default_pretty=False
    )
    main(parser, process_raw_genesis)
