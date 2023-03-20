package keeper

import (
	"github.com/classic-terra/core/x/treasury/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var burnTaxExcemptionAddressList = []string{
	"terra10atxpzafqfjy58z0dvugmd9zf63fycr6uvwhjm",
	"terra1jrq7xa63a4qgpdgtj70k8yz5p32ps9r7mlj3yr",
	"terra15s66unmdcpknuxxldd7fsr44skme966tdckq8c",
	"terra1u0p7xuwlg0zsqgntagdkyjyumsegd8agzhug99",
	"terra1fax8l6srhew5tu2mavmu83js3v7vsqf9yr4fv7",
	"terra132wegs0kf9q65t9gsm3g2y06l98l2k4treepkq",
	"terra1l89hprfccqxgzzypjzy3fnp7vqpnkqg5vvqgjc",
	"terra1ns7lfvrxzter4d2yl9tschdwntcxa25vtsvd8a",
	"terra1vuvju6la7pj6t8d8zsx4g8ea85k2cg5u62cdhl",
	"terra1lzdux37s4anmakvg7pahzh03zlf43uveq83wh2",
	"terra1ky3qcf7v45n6hwfmkm05acwczvlq8ahnq778wf",
	"terra17m8tkde0mav43ckeehp537rsz4usqx5jayhf08",
	"terra1urj8va62jeygra7y3a03xeex49mjddh3eul0qa",
	"terra10wyptw59xc52l86pg86sy0xcm3nm5wg6a3cf7l",
	"terra1sujaqwaw7ls9fh6a4x7n06nv7fxx5xexwlnrkf",
	"terra1qg59nhvag222kp6fyzxt83l4sw02huymqnklww",
	"terra1dxxnwxlpjjkl959v5xrghx0dtvut60eef6vcch",
	"terra1y246m036et7vu69nsg4kapelj0tywe8vsmp34d",
	"terra1j39c9sjr0zpjnrfjtthuua0euecv7txavxvq36",
	"terra1t0jthtq9zhm4ldtvs9epp02zp23f355wu6zrzq",
	"terra12dxclvqrgt7w3s7gtwpdkxgymexv8stgqcr0yu",
	"terra1az3dsad74pwhylrrexnn5qylzj783uyww2s7xz",
	"terra1ttq26dq4egr5exmhd6gezerrxhlutx9u90uncn",
	"terra13e9670yuvfs06hctt9pmgjnz0yw28p0wgnhrqn",
	"terra1skmktm537pfaycgu9jx4fqryjt6pf77ycpesw0",
	"terra14q8cazgt58y2xkd26mlukemwth0cnvfqmgz2qk",
	"terra163vzxz9wwy320ccwy73qe6h33yzg2yhyvv5nsf",
	"terra1kj43wfnvrgc2ep94dgmwvnzv8vnkkxrxmrnhkp",
	"terra1gu6re549pn0mdpshtv75t3xugn347jghlhul73",
	"terra1gft3qujlq04yza3s2r238mql2yn3xxqepzt2up",
	"terra174pe7qe7g867spzdfs5f4rf9fuwmm42zf4hykf",
	"terra1ju68sg6k39t385sa0fazqvjgh6m6gkmsmp4lln",
	"terra1dlh7k4hcnsrvlfuzhdzx3ctynj7s8dde9zmdyd",
	"terra18wcdhpzpteharlkks5n6k7ent0mjyftvcpm6ee",
	"terra1xmkwsauuk3kafua9k23hrkfr76gxmwdfq5c09d",
	"terra1t957gces65xd6p8g4cuqnyd0sy5tzku59njydd",
	"terra1s4rd0y5e4gasf0krdm2w8sjhsmh030m74f2x9v",
	"terra15jya6ugxp65y80y5h82k4gv90pd7acv58xp6jj",
	"terra14yqy9warjkxyecda5kf5a68qlknf4ve4sh7sa6",
	"terra1yxras4z0fs9ugsg2hew9334k65uzejwcslyx0y",
	"terra1p0vl4s4gp46vy6dm352s2fgtw6hccypph7zc3u",
	"terra1hhj92twle9x8rjkr3yffujexsy5ldexak5rglz",
	"terra18vnrzlzm2c4xfsx382pj2xndqtt00rvhu24sqe",
	"terra1ncjg4a59x2pgvqy9qjyqprlj8lrwshm0wleht5",
	"terra19l7hzwazq5j0dykfldcwrk2927xwcjd0kt0vt9",
	"terra1frh79vmtur5fmrghz6gfjvfhpa3u2c0uemv4af",
}

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	m.keeper.SetBurnSplitRate(ctx, types.DefaultBurnTaxSplit)

	for _, address := range burnTaxExcemptionAddressList {
		m.keeper.AddBurnTaxExemptionAddress(ctx, address)
	}

	return nil
}

// Migrate2to3 migrates from version 2 to 3.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	m.keeper.SetMinInitialDepositRatio(ctx, types.DefaultMinInitialDepositRatio)

	return nil
}
