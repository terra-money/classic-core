
#!/usr/bin/env python3

import operator
import argparse
import json
import sys
import dateutil.parser
import math

DECIMAL_UNIT = 3

def init_default_argument_parser(prog_desc, default_chain_id, default_genesis_time):
    parser = argparse.ArgumentParser(description=prog_desc)
    parser.add_argument(
        'exported_genesis',
        help='exported genesis.json file',
        type=argparse.FileType('r'), default=sys.stdin,
    )
    parser.add_argument('--chain-id', type=str, default=default_chain_id)
    parser.add_argument('--genesis-time', type=str, default=default_genesis_time)
    return parser

# for coins sort
def coin_sort(val):
  return val['denom']

def add_coin(coins, coin):
    for c in coins:
        if c['denom'] == coin['denom']:
            c['amount'] = str(int(c['amount']) + int(coin['amount']))
            return

    coins.append(coin)
    coins.sort(key=coin_sort)
  

def main(argument_parser, process_genesis_func):
    args = argument_parser.parse_args()
    if args.chain_id.strip() == '':
      sys.exit('chain-id required')

    genesis = json.loads(args.exported_genesis.read())

    genesis = process_genesis_func(genesis=genesis, parsed_args=args,)

    # update terra1n2kzv00yjanjpjplqtwucug45lurr8tzgrvj2p => terra1pnnruerze80znqdcl5ka6y3ntpzyu2e4j4yeev
    # update terra1f0gxt604sn3py38u4pecf2ha3x66640cvndpv7 => terra1pln09kv8s7k3ssw0p6ymdwjn7aya3h6mj5mp2x
    raw_genesis = json.dumps(genesis, indent=4, sort_keys=True)
    raw_genesis = raw_genesis.replace(
        'terra1n2kzv00yjanjpjplqtwucug45lurr8tzgrvj2p', 
        'terra1pnnruerze80znqdcl5ka6y3ntpzyu2e4j4yeev')
    raw_genesis = raw_genesis.replace(
        'terra1f0gxt604sn3py38u4pecf2ha3x66640cvndpv7',
        'terra1pln09kv8s7k3ssw0p6ymdwjn7aya3h6mj5mp2x')

    print(raw_genesis)


def create_module_account(name, address, coins, permissions):
    return {
        'address': address,
        'coins': coins,
        'sequence_number': '0',
        'account_number': '0',
        'original_vesting': [],
        'delegated_free': [],
        'delegated_vesting': [],
        'start_time': '0',
        'end_time': '0',
        'module_name': name,
        'module_permissions': permissions,
        'vesting_schedules': []
    }

def process_raw_genesis(genesis, parsed_args):

    # Delete bugdet
    del genesis['app_state']['budget']

    bondedAmt = 0
    notBondedAmt = 0
    for val in genesis['app_state']['staking']['validators']:
        if val['status'] == 2:
            bondedAmt += int(val['tokens'])
        elif val['status'] == 0 or val['status'] == 1:
            notBondedAmt += int(val['tokens'])
        else:
            sys.exit('Invalid validator status')
    
    if genesis['app_state']['staking']['unbonding_delegations']:
        for ubd in genesis['app_state']['staking']['unbonding_delegations']:
            for entry in ubd['entries']:
                notBondedAmt += int(entry['balance'])

    newAccounts = []
    # Change old genesis accounts to new format genesis account
    for acc in genesis['app_state']['accounts']:
        newAcc = {
            'address': acc['address'],
            'coins': acc['coins'],
            'sequence_number': acc['sequence_number'],
            'account_number': '0', # ignored set by the account keeper during InitGenesis
            'original_vesting': acc['original_vesting'],
            'delegated_free': acc['delegated_free'],
            'delegated_vesting': acc['delegated_vesting'],
            'start_time': '0',
            'end_time': '0',
            'module_name': '',
            'module_permissions': [],
            'vesting_schedules': acc['lazy_vesting_schedules']
        }

        if acc['address'] == 'terra1fs7mmpducjf25j70sk3sz6k5phz2fllmyr5gwz':
            update_vesting_schedule(newAcc)
        if acc['address'] == 'terra1dp0taj85ruc299rkdvzp4z5pfg6z6swaed74e6':
            # (Foundation) sub 978,260 LUNA to give new vesting account
            for coin in newAcc['coins']:
              if coin['denom'] == 'uluna':
                coin['amount'] = str(int(coin['amount']) - 978260000000)
        if acc['address'] == 'terra1nl2vrxr0qzzy4pd9m2mw0q0tvwcxe2mg8shaad':
            update_seed_to_private_vesting_schedule(newAcc)
        if acc['address'] == 'terra10y5usrnwk2ltddm5kenhznl5uj6w3yfga5al4a':
            # 978,260 LUNA 1M, 2M, 3M, 12M 10% 10% 10% 70%
            coin = {
              'denom': 'uluna',
              'amount': '978260000000'
            }

            add_coin(newAcc['coins'], coin)
            newAcc['original_vesting'] = [coin]
            update_seed2_vesting_schedule(newAcc)

        newAccounts.append(newAcc)

    # Create fee collector account
    newAccounts.append(
        create_module_account(
            'fee_collector', 
            'terra17xpfvakm2amg962yls6f84z3kell8c5lkaeqfa', 
            genesis['app_state']['auth']['collected_fees'], 
            ['basic']
        )
    )

    # Create empty gov account with burner permissions
    newAccounts.append(
        create_module_account(
            'gov', 
            'terra10d07y265gmmuvt4z0w9aw880jnsr700juxf95n', 
            [], 
            ['burner']
        )
    )

    # Create empty distribution account
    # We assume genesis is exported for zero-height without any remaining rewards
    # only cares for community-pool
    communityPoolCoins = []
    for coin in genesis['app_state']['distr']['fee_pool']['community_pool']:
        communityPoolCoins.append({
            'denom': coin['denom'],
            'amount': str(int(float(coin['amount'])))
        })

    newAccounts.append(
        create_module_account(
            'distribution', 
            'terra1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8pm7utl', 
            communityPoolCoins, 
            ['basic']
        )
    )

    # Create bonded_tokens_pool account with burner & staking permissions
    newAccounts.append(
        create_module_account(
            'bonded_tokens_pool', 
            'terra1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3nln0mh', 
            [{'amount': str(bondedAmt), 'denom': 'uluna'}], 
            ['burner', 'staking']
        )
    )

    # Create not_bonded_tokens_pool account with burner & staking permissions
    newAccounts.append(
        create_module_account(
            'not_bonded_tokens_pool', 
            'terra1tygms3xhhs3yv487phx3dw4a95jn7t7l8l07dr', 
            [{'amount': str(notBondedAmt), 'denom': 'uluna'}], 
            ['burner', 'staking']
        )
    )

    # Create empty oracle account
    newAccounts.append(
        create_module_account(
            'oracle', 
            'terra1jgp27m8fykex4e4jtt0l7ze8q528ux2lh4zh0f', 
            [], 
            ['basic']
        )
    )

    # Create empty market account
    newAccounts.append(
        create_module_account(
            'market', 
            'terra1untf85jwv3kt0puyyc39myxjvplagr3wstgs5s', 
            [], 
            ['minter', 'burner']
        )
    )

    # Create empty treasury account
    newAccounts.append(
        create_module_account(
            'treasury', 
            'terra1vmafl8f3s6uuzwnxkqz0eza47v6ecn0t0yeca7', 
            [], 
            ['minter']
        )
    )

    # Migrate accounts 
    genesis['app_state']['accounts'] = newAccounts

    # Migrate auth module 
    genesis['app_state']['auth'] = {
        'params': genesis['app_state']['auth']['params']
    }

    # Add gov module genesis state
    genesis['app_state']['gov'] = {
        'deposits': [],
        'proposals': [],
        'votes': [],
        'starting_proposal_id': '1',
        'deposit_params': {
            'max_deposit_period': '1209600000000000',
            'min_deposit': [{
                'amount': '512000000',
                'denom': 'uluna'
            }]
        },
        'voting_params': {
            'voting_period': '1209600000000000'
        },
        'tally_params': {
            'quorum': '0.400000000000000000',
            'threshold': '0.500000000000000000',
            'veto': '0.334000000000000000'
        }
    }

    # Increase slashing panelty for security improvement
    genesis['app_state']['slashing']['params']['slash_fraction_double_sign'] = '0.05'

    # Move genesis state key from distr to distribution
    genesis['app_state']['distribution'] = genesis['app_state']['distr']
    del genesis['app_state']['distr']
    
    # Migrate distribution module
    for vse in genesis['app_state']['distribution']['validator_slash_events']:
        vse['period'] = vse['validator_slash_event']['validator_period']

    # Migrate staking module
    for val in genesis['app_state']['staking']['validators']:
        val['commission'] = {
            'commission_rates': {
                'rate': val['commission']['rate'],
                'max_rate': val['commission']['max_rate'],
                'max_change_rate': val['commission']['max_change_rate']
            },
            'update_time': val['commission']['update_time']
        }

    del genesis['app_state']['staking']['pool']

    # Add supply module genesis state
    genesis['app_state']['supply'] = {
        'supply': []
    }

    # Migrate market module
    genesis['app_state']['market'] = {
        'terra_pool_delta': '0',
        'params': {
            'base_pool': '250000000000',          # 250,000 sdr = 250,000,000,000 usdr
            'pool_recovery_period': '14400',      # blocks per day
            'min_spread': '0.020000000000000000', # 2%
            'tobin_tax': '0.002500000000000000',  # 0.25%
            'illiquid_tobin_tax_list': [{
                'denom': 'umnt',
                'tax_rate': '0.02'                # 2%
            }]
        }
    }

    # Migrate oracle module
    genesis['app_state']['oracle'] = {
        'feeder_delegations': genesis['app_state']['oracle']['feeder_delegations'],
        'exchange_rates': {},
        'exchange_rate_prevotes': [],
        'exchange_rate_votes': [],
        'miss_counters': {},
        'params': {
            'vote_period': '5', # blocks per half-minute
            'vote_threshold': genesis['app_state']['oracle']['params']['vote_threshold'],
            'whitelist': ['ukrw', 'usdr', 'uusd', 'umnt'],
            'reward_band': '0.07',                     # 7%
            'reward_distribution_window': '5256000',   # num of blocks per a year
            'slash_window': '432000',                  # num of blocks per a month
            'slash_fraction': '0.0001',                # 0.01%
            'min_valid_per_window': '0.05'             # 5%
        }
    }

    # Migrate treasury module
    genesis['app_state']['treasury'] = {
        'params': genesis['app_state']['treasury']['params'],
        'tax_rate': '0.005', # 0.5%
        'reward_weight': genesis['app_state']['treasury']['reward_weight'],
        'tax_caps': {
            'ukrw': '1612000000',    # 1612 KRW
            'uusd': '1370000',       # 1.37 USD
            'umnt': '3715000000',    # 3715 MNT
            'usdr': '1000000',       # 1 SDR
        },
        'tax_proceed': [],
        'epoch_initial_issuance': [],
        'TRs': [],
        'SRs': [],
        'TSLs': []
    }

    genesis['app_state']['treasury']['params']['tax_policy']['rate_min'] = '0.001' # 0.1%
    genesis['app_state']['treasury']['params']['reward_policy']['rate_max'] = '0.5' # 50%

    # TODO - should we need to add window_probation?
    # del genesis['app_state']['treasury']['params']['window_probation']
    genesis['app_state']['treasury']['params']['window_probation'] = '0'
    
    # Set new chain ID and genesis start time
    genesis['chain_id'] = parsed_args.chain_id.strip()
    genesis['genesis_time'] = parsed_args.genesis_time

    return genesis

def update_vesting_schedule(account):
    initial_genesis_time = '2019-04-24T06:00:00.000000Z'
    
    # parse genesis date
    genesis_date = dateutil.parser.parse(initial_genesis_time)
    vesting_schedules = []

    # Luna Schedule Update
    luna_vesting_schedule = {
        'denom': 'uluna',
        'schedules': [
            {
                'start_time': str(get_time_after_n_month(genesis_date, 4)),
                'end_time': str(get_time_after_n_month(genesis_date, 5)),
                'ratio': '0.100000000000000000',
            },
            {
                'start_time': str(get_time_after_n_month(genesis_date, 5)),
                'end_time': str(get_time_after_n_month(genesis_date, 6)),
                'ratio': '0.100000000000000000',
            },
            {
                'start_time': str(get_time_after_n_month(genesis_date, 6)),
                'end_time': str(get_time_after_n_month(genesis_date, 7)),
                'ratio': '0.100000000000000000',
            },
            {
                'start_time': str(get_time_after_n_month(genesis_date, 24)),
                'end_time': str(get_time_after_n_month(genesis_date, 25)),
                'ratio': '0.700000000000000000',
            }
        ]
    }

    # Terra Schedule Update
    terra_schedules = []
    cumulated_ratio = 0
    for i in range(17):
        ratio = correct_decimal(1/18)
        cumulated_ratio += ratio
        terra_schedules.append({
            'start_time': str(get_time_after_n_month(genesis_date, 4 + i)),
            'end_time': str(get_time_after_n_month(genesis_date, 5 + i)),
            'ratio': str(ratio),
        })
    

    terra_schedules.append({
        'start_time': str(get_time_after_n_month(genesis_date, 4 + 17)),
        'end_time': str(get_time_after_n_month(genesis_date, 5 + 17)),
        'ratio': str(1 - cumulated_ratio),
    })

    terra_vesting_schedule = {
        'denom': 'usdr',
        'schedules': terra_schedules
    }

    vesting_schedules.append(luna_vesting_schedule)
    vesting_schedules.append(terra_vesting_schedule)
    account['vesting_schedules'] = vesting_schedules

def update_seed_to_private_vesting_schedule(account):
    initial_genesis_time = '2019-04-24T06:00:00.000000Z'
    
    # parse genesis date
    genesis_date = dateutil.parser.parse(initial_genesis_time)
    vesting_schedules = []

    # Luna Schedule Update
    luna_vesting_schedule = {
        'denom': 'uluna',
        'schedules': [
            {
                'start_time': str(get_time_after_n_month(genesis_date, 4)),
                'end_time': str(get_time_after_n_month(genesis_date, 5)),
                'ratio': '0.166000000000000000',
            },
            {
                'start_time': str(get_time_after_n_month(genesis_date, 5)),
                'end_time': str(get_time_after_n_month(genesis_date, 6)),
                'ratio': '0.166000000000000000',
            },
            {
                'start_time': str(get_time_after_n_month(genesis_date, 6)),
                'end_time': str(get_time_after_n_month(genesis_date, 7)),
                'ratio': '0.166000000000000000',
            },
            {
                'start_time': str(get_time_after_n_month(genesis_date, 7)),
                'end_time': str(get_time_after_n_month(genesis_date, 8)),
                'ratio': '0.166000000000000000',
            },
            {
                'start_time': str(get_time_after_n_month(genesis_date, 8)),
                'end_time': str(get_time_after_n_month(genesis_date, 9)),
                'ratio': '0.166000000000000000',
            },
            {
                'start_time': str(get_time_after_n_month(genesis_date, 9)),
                'end_time': str(get_time_after_n_month(genesis_date, 10)),
                'ratio': '0.17000000000000000',
            }
        ]
    }


    # Terra vesting has no need to be updated
    terra_vesting_schedule = {
        'denom': 'usdr'
    }

    # Find origin terra vesting schedule and use it to new vesting schedule
    for vesting_schedule in account['vesting_schedules']:
        if vesting_schedule['denom'] == 'usdr':
            terra_vesting_schedule['schedules'] = vesting_schedule['schedules']
            break

    vesting_schedules.append(luna_vesting_schedule)
    vesting_schedules.append(terra_vesting_schedule)
    account['vesting_schedules'] = vesting_schedules

def update_seed2_vesting_schedule(account):
    initial_genesis_time = '2019-12-13T00:00:00.000000Z'
    
    # parse genesis date
    genesis_date = dateutil.parser.parse(initial_genesis_time)
    vesting_schedules = []

    # Luna Schedule Update
    luna_vesting_schedule = {
        'denom': 'uluna',
        'schedules': [
            {
                'start_time': str(get_time_after_n_month(genesis_date, 1)),
                'end_time': str(get_time_after_n_month(genesis_date, 2)),
                'ratio': '0.100000000000000000',
            },
            {
                'start_time': str(get_time_after_n_month(genesis_date, 2)),
                'end_time': str(get_time_after_n_month(genesis_date, 3)),
                'ratio': '0.100000000000000000',
            },
            {
                'start_time': str(get_time_after_n_month(genesis_date, 3)),
                'end_time': str(get_time_after_n_month(genesis_date, 4)),
                'ratio': '0.100000000000000000',
            },
            {
                'start_time': str(get_time_after_n_month(genesis_date, 12)),
                'end_time': str(get_time_after_n_month(genesis_date, 13)),
                'ratio': '0.700000000000000000',
            }
        ]
    }

    vesting_schedules.append(luna_vesting_schedule)
    account['vesting_schedules'] = vesting_schedules


def get_time_after_n_month(start_date, n):
    year = start_date.year
    month = start_date.month+n
    while month > 12:
        year += 1
        month -= 12

    return int(start_date.replace(year=year, month=month).timestamp())

def correct_decimal(float_num):
    return int(float_num * (10**DECIMAL_UNIT)) / (10**DECIMAL_UNIT)

if __name__ == '__main__':
    parser = init_default_argument_parser(
        prog_desc='Convert genesis.json for columbus-3',
        default_chain_id='columbus-3',
        default_genesis_time='2019-12-13T15:00:00Z',
    )
    main(parser, process_raw_genesis)
