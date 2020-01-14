# Stability

The basic principle behind Terra’s stability has always been the same: adjust its supply in accordance to changes in demand to keep its price stable \(say at 1SDR\), and do this in a decentralized and algorithmic way. More precisely:

* When the price of Terra rises above the target the protocol needs to expand its supply
* When its price falls below the target the protocol needs to contract its supply

The latter is the difficult part of course — printing money is easy when you hold the printing press.

Contracting and expanding Terra’s supply simply means buying or selling Terra at its target price respectively. The protocol uses Luna to facilitate those exchanges, essentially acting as a market maker for Terra:

* To buy 1 Terra, the protocol mints and sells Luna worth 1 SDR
* By selling 1 Terra, the protocol earns Luna worth 1 SDR

Insofar as the protocol can maintain both sides of this contract, the price of Terra ought to remain pegged at 1 SDR.

What is simply happening here is that the supply of Luna is absorbing volatility in the demand for Terra, meaning that the supply of Luna decreases or increases in accordance to changes in Terra demand. This of course affects miners by either burning or diluting mining power. Mining power becomes more or less valuable in the process and so in effect miners are the ones absorbing Terra’s immediate volatility.

Miners are essential for the security of the network and the stability of Terra. The protocol is thus designed to maintain stable demand for mining in the long-term. This is achieved by shielding miners from the volatility of Terra’s economic cycles — essentially rewarding them more when demand for Terra decreases and less when demand for Terra increases to smooth out the ups and downs. Miners are rewarded in two ways:

* **Transaction Fees**: Every Terra transaction pays a small fee to miners. Fees vary over time to smooth out mining rewards. For example, if mining rewards are declining, an increase in fees can reverse that trend and vice versa.
* **Seigniorage**: Whenever the protocol mints Terra and earns Luna in return it generates what we call “seigniorage” — the profit from minting currency. A portion of earned Luna is burned as a way to reward miners. As with transaction fees, the portion of seigniorage that rewards miners is adjusted over time to smooth out fluctuations. Remaining seigniorage is allocated to the Treasury for stimulating growth — we explain how this works in the following section.

To sum up, the protocol makes it a priority to shield miners from excessive volatility in Terra demand by making mining rewards countercyclical — controlling them to oppose Terra’s economic cycles.

