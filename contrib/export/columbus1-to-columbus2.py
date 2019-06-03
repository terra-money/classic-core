#!/usr/bin/env python3

import lib
import sys
import dateutil.parser

DECIMAL_UNIT = 3
SECONDS_PER_MINUTE = 60
SECONDS_PER_HOURS = SECONDS_PER_MINUTE * 60
SECONDS_PER_DAY = SECONDS_PER_HOURS * 24
SECONDS_PER_MONTH = SECONDS_PER_HOURS * 30

def process_raw_genesis(genesis, vesting_info, parsed_args):

    # parse genesis date
    genesis_date = dateutil.parser.parse(genesis['genesis_time'])

    # modify i-4 type accounts vesting schedule
    accounts = genesis['app_state']['accounts']

    for account in accounts:
        for row in vesting_info:
            if len(row) != 3:
                sys.exit('invalid cvs format')

            if account['address'] == row[0]:
                change_vesting_schedule(account, row[1], row[2], genesis_date)

        # Change wallet address of translink capital
        if account['address'] == 'terra1d4v2k6h7ltc4mpkk9pgtp277yp2flte7296ypq':
            account['address'] = 'terra1vulcsvfla6sej0p5y0n99w4ldke5tte3rsukyk'


    # Set new chain ID and genesis start time
    genesis['chain_id'] = parsed_args.chain_id.strip()
    genesis['genesis_time'] = parsed_args.start_time

    return genesis

def change_vesting_schedule(account, denom, amount, genesis_date):

    ratio = False
    original_vesting = account['original_vesting']
    for vesting in original_vesting:
        if vesting['denom'] == denom:
            ratio = int(amount) / int(vesting['amount'][:-6])
            break

    if ratio == False:
        sys.exit('invalid account or denom given')

    vesting_schedules = account['vesting_schedules']
    for vesting_schedule in vesting_schedules:
        
        if vesting_schedule['denom'] == denom:

            schedules =   vesting_schedule['schedules']
            for schedule in schedules:
                if schedule['cliff'] == str(get_time_after_n_month(genesis_date, 4)) or schedule['cliff'] == str(get_time_after_n_month(genesis_date, 5)):
                    schedule['ratio'] = str(float(schedule['ratio']) - (ratio / 2))

            for i in range(4, 10):
                cliff = get_time_after_n_month(genesis_date, i)
                single_ratio = correct_decimal(ratio / 6)

                found = False
                for schedule in schedules:
                    if schedule['cliff'] == str(cliff):
                        schedule['ratio'] = str(correct_decimal(float(schedule['ratio']) + single_ratio))
                        found = True
                        break
                
                if not found:
                    schedules.append({
                        'cliff': str(cliff),
                        'ratio': str(single_ratio)
                    })
            
            schedules.sort(key=custom_sort)

            ratio_sum = 0
            for schedule in schedules:
                ratio_sum += float(schedule['ratio'])

            last_index = len(schedules) - 1
            schedules[last_index]['ratio'] = str(round(float(schedules[last_index]['ratio']) - (ratio_sum - 1), DECIMAL_UNIT))

            invariant_check(schedules)

def custom_sort(schedule):
    return schedule['cliff']

def correct_decimal(float_num):
    return int(float_num * (10**DECIMAL_UNIT)) / (10**DECIMAL_UNIT)

def invariant_check(schedules):
    ratio_sum = 0
    for schedule in schedules:
        ratio_sum += float(schedule['ratio'])

    if ratio_sum != 1:
        sys.exit('%f, invariant failed' % ratio_sum)


def get_time_after_n_month(start_date, n):
    target_month = start_date.month+n
    year = target_month > 12 and start_date.year + int(target_month / 12) or start_date.year
    month = target_month - 12 * (year - start_date.year)
    return int(start_date.replace(year=year, month=month).timestamp())

if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Convert genesis.json from columbus-1 to columbus-2',
        default_chain_id='columbus-2',
        default_start_time='2019-05-20T05:00:00Z',
    )
    lib.main(parser, process_raw_genesis)
