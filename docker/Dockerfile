ARG base_image=terraclassic.terrad-binary
ARG version=latest

FROM ${base_image}:${version}

ARG chainid=columbus-5
ARG enable_lcd=true
ARG enable_unsafe_cors=false
ARG cors_allowed_origins=[]
ARG cors_allowed_methods=["HEAD", "GET", "POST"]
ARG cors_allowed_headers=["Origin", "Accept", "Content-Type", "X-Requested-With", "X-Server-Time"]
ARG validator_keyname=local
ARG validator_passphrase=12345678
ARG validator_mnenomic="torch swamp cancel lift never october child harsh rib aspect luxury word peanut length bamboo hawk material vehicle glue above west random sketch author"
ARG validator_amount=1uluna
ARG validator_commission_rate=0.2
ARG validator_commission_rate_max=1
ARG validator_commission_rate_max_change=0.01
ARG validator_min_self_delegation=1

ENV CHAINID ${chainid}
ENV ENABLE_LCD ${enable_lcd}
ENV ENABLE_UNSAFE_CORS ${enable_unsafe_cors}
ENV CORS_ALLOWED_ORIGINS ${cors_allowed_origins}
ENV CORS_ALLOWED_METHODS ${cors_allowed_methods}
ENV CORS_ALLOWED_HEADERS ${cors_allowed_headers}
ENV VALIDATOR_KEYNAME ${validator_keyname}
ENV VALIDATOR_PASSPHRASE ${validator_passphrase}
ENV VALIDATOR_MNENOMIC ${validator_mnenomic}
ENV VALIDATOR_AMOUNT ${validator_amount}
ENV VALIDATOR_COMMISSION_RATE ${validator_commission_rate}
ENV VALIDATOR_COMMISSION_RATE_MAX ${validator_commission_rate_max}
ENV VALIDATOR_COMMISSION_RATE_MAX_CHANGE ${validator_commission_rate_max_change}
ENV VALIDATOR_MIN_SELF_DELEGATION ${validator_min_self_delegation}

# Moniker will be updated by entrypoint.
RUN terrad init --chain-id $chainid moniker

# Backup for templating
RUN mv ~/.terra/config/config.toml ~/config.toml
RUN mv ~/.terra/config/app.toml ~/app.toml

RUN if [ "$chainid" = "columbus-5" ] ; then wget -O ~/.terra/config/genesis.json https://columbus-genesis.s3.ap-northeast-1.amazonaws.com/columbus-5-genesis.json; fi; \
    if [ "$chainid" = "columbus-5" ] ; then wget -O ~/.terra/config/addrbook.json https://networks.mcontrol.ml/columbus/addrbook.json; fi; \
    if [ "$chainid" = "rebel-1" ]    ; then wget -O ~/.terra/config/genesis.json https://raw.githubusercontent.com/terra-rebels/classic-testnet/master/rebel-1/genesis.json; fi; \
    if [ "$chainid" = "rebel-1" ]    ; then wget -O ~/.terra/config/addrbook.json https://raw.githubusercontent.com/terra-rebels/classic-testnet/master/rebel-1/addrbook.json; fi; \
    if [ "$chainid" = "rebel-2" ]    ; then wget -O ~/.terra/config/genesis.json https://raw.githubusercontent.com/terra-rebels/classic-testnet/master/rebel-2/genesis.json; fi; \
    if [ "$chainid" = "rebel-2" ]    ; then wget -O ~/.terra/config/addrbook.json https://raw.githubusercontent.com/terra-rebels/classic-testnet/master/rebel-2/addrbook.json; fi;
    
COPY ./entrypoint.sh /entrypoint.sh
COPY ./keys-add.sh /keys-add.sh
COPY ./create-validator.sh /create-validator.sh

CMD ["/entrypoint.sh"]