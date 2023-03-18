ARG base_image=classic-terra/core
ARG version=latest

FROM ${base_image}:${version}

ARG chainid=test
ARG new_network=false
ARG enable_lcd=true
ARG validator_keyname=local
ARG validator_mnenomic="torch swamp cancel lift never october child harsh rib aspect luxury word peanut length bamboo hawk material vehicle glue above west random sketch author"
ARG validator_amount=1uluna
ARG validator_commission_rate=0.2
ARG validator_commission_rate_max=1
ARG validator_commission_rate_max_change=0.01
ARG validator_min_self_delegation=1

ENV CHAINID ${chainid}
ENV NEW_NETWORK ${new_network}
ENV ENABLE_LCD ${enable_lcd}
ENV VALIDATOR_KEYNAME ${validator_keyname}
ENV VALIDATOR_MNENOMIC ${validator_mnenomic}
ENV VALIDATOR_AMOUNT ${validator_amount}
ENV VALIDATOR_COMMISSION_RATE ${validator_commission_rate}
ENV VALIDATOR_COMMISSION_RATE_MAX ${validator_commission_rate_max}
ENV VALIDATOR_COMMISSION_RATE_MAX_CHANGE ${validator_commission_rate_max_change}
ENV VALIDATOR_MIN_SELF_DELEGATION ${validator_min_self_delegation}

COPY ./entrypoint.sh /entrypoint.sh
COPY ./keys-add.sh /keys-add.sh
COPY ./create-validator.sh /create-validator.sh
COPY ./test-node-setup.sh /test-node-setup.sh
COPY ./keys.json /keys.json

CMD ["/entrypoint.sh"]