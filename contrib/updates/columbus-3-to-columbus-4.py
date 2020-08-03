#!/usr/bin/env python3

import argparse
import json
import sys


def init_default_argument_parser(prog_desc, default_chain_id, default_genesis_time):
    parser = argparse.ArgumentParser(description=prog_desc)
    parser.add_argument(
        'exported_genesis',
        help='exported genesis.json file',
        type=argparse.FileType('r'), default=sys.stdin,
    )
    parser.add_argument('--chain-id', type=str, default=default_chain_id)
    parser.add_argument('--genesis-time', type=str,
                        default=default_genesis_time)
    return parser


def main(argument_parser, process_genesis_func):
    args = argument_parser.parse_args()
    if args.chain_id.strip() == '':
        sys.exit('chain-id required')

    genesis = json.loads(args.exported_genesis.read())
    genesis = process_genesis_func(genesis=genesis, parsed_args=args,)

    raw_genesis = json.dumps(genesis, indent=4, sort_keys=True)
    print(raw_genesis)


def process_raw_genesis(genesis, parsed_args):

    # Auth Module Migration
    accs = []
    for acc in genesis['app_state']['accounts']:
        accType = 'core/Account'
        if acc['module_name']:
            # Module Acc
            accType = 'supply/ModuleAccount'
            del acc['delegated_free']
            del acc['delegated_vesting']
            del acc['end_time']
            del acc['original_vesting']
            del acc['start_time']
            del acc['vesting_schedules']
        elif acc['vesting_schedules']:
            # Vesting Acc
            accType = 'core/LazyGradedVestingAccount'
            del acc['module_name']
            del acc['module_permissions']
            del acc['start_time']
        else:
            # Normal Acc
            del acc['delegated_free']
            del acc['delegated_vesting']
            del acc['end_time']
            del acc['original_vesting']
            del acc['start_time']
            del acc['vesting_schedules']
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

    # Migrate Tendermint Consensus Param
    genesis['consensus_params']['evidence'] = {
        'max_age_num_blocks': genesis['consensus_params']['evidence']['max_age'],
        # 2 days, should we increase it; https://github.com/tendermint/tendermint/issues/2565
        'max_age_duration': '172800000000000',
    }

    # Set new chain ID and genesis start time
    genesis['chain_id'] = parsed_args.chain_id.strip()
    genesis['genesis_time'] = parsed_args.genesis_time

    return genesis


if __name__ == '__main__':
    parser = init_default_argument_parser(
        prog_desc='Convert genesis.json for columbus-4',
        default_chain_id='columbus-4',
        default_genesis_time='2020-08-01T15:00:00Z',
    )
    main(parser, process_raw_genesis)
