package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

// Constants containing the genesis allocation of built-in genesis blocks.
// Copied from https://github.com/celo-org/celo-blockchain/blob/d4ac28c04d5b94855e1e72c9b4436fa3c73e16ad/core/genesis_alloc.go#L19

const baklavaAllocJSON = `{
    "fCf982bb4015852e706100B14E21f947a5Bb718E": {
      "balance": "200000000000000000000000000"
    },
    "0xd71fea6b92d3f21f659152589223385a7329bb11": {
      "balance": "1000000000000000000000"
    },
    "0x1e477fc9b6a49a561343cd16b2c541930f5da7d2": {
      "balance": "1000000000000000000000"
    },
    "0x460b3f8d3c203363bb65b1a18d89d4ffb6b0c981": {
      "balance": "1000000000000000000000"
    },
    "0x3b522230c454ca9720665d66e6335a72327291e8": {
      "balance": "1000000000000000000000"
    },
    "0x0AFe167600a5542d10912f4A07DFc4EEe0769672": {
      "balance": "1000000000000000000000"
    },
    "0x412ebe7859e9aa71ff5ce4038596f6878c359c96": {
      "balance": "1000000000000000000000"
    },
    "0xbbfe73df8b346b3261b19ac91235888aba36d68c": {
      "balance": "1000000000000000000000"
    },
    "0x02b1d1bea682fcab4448c0820f5db409cce4f702": {
      "balance": "1000000000000000000000"
    },
    "0xe90f891710f625f18ecbf1e02efb4fd1ab236a10": {
      "balance": "1000000000000000000000"
    },
    "0x28c52c722df87ed11c5d7665e585e84aa93d7964": {
      "balance": "1000000000000000000000"
    },
    "0Cc59Ed03B3e763c02d54D695FFE353055f1502D": {
      "balance": "103010030000000000000000000"
    },
    "3F5084d3D4692cf19b0C98A9b22De614e49e1470": {
      "balance": "10011000000000000000000"
    },
    "EF0186B8eDA17BE7D1230eeB8389fA85e157E1fb": {
      "balance": "10011000000000000000000"
    },
    "edDdb60EF5E90Fb09707246DF193a55Df3564c9d": {
      "balance": "10011000000000000000000"
    },
    "d5e454462b3Fd98b85640977D7a5C783CA162228": {
      "balance": "10011000000000000000000"
    },
    "a4f1bad7996f346c3E90b90b60a1Ca8B67B51E4B": {
      "balance": "10011000000000000000000"
    },
    "5B991Cc1Da0b6D54F8befa9De701d8BC85C92324": {
      "balance": "10011000000000000000000"
    },
    "6dfdAa51D146eCff3B97614EF05629EA83F4997E": {
      "balance": "10011000000000000000000"
    },
    "D2b16050810600296c9580D947E9D919D0c332ed": {
      "balance": "10011000000000000000000"
    },
    "Fe144D67068737628efFb701207B3eB30eF93C69": {
      "balance": "10011000000000000000000"
    },
    "82E64996B355625efeAaD12120710706275b5b9A": {
      "balance": "10011000000000000000000"
    },
    "241752a3f65890F4AC3eAeC518fF94567954e7b5": {
      "balance": "10011000000000000000000"
    },
    "1bdDeaF571d5da96ce6a127fEb3CADaDB531f433": {
      "balance": "10011000000000000000000"
    },
    "F86345e9c9b39aB1cbE82d7aD35854f905B8B835": {
      "balance": "10011000000000000000000"
    },
    "5c3512b1697302c497B861CBfDA158f8a3c5122C": {
      "balance": "10011000000000000000000"
    },
    "a02A692d70Fd9A5269397C044aEBDf1085ba090f": {
      "balance": "10011000000000000000000"
    },
    "aC91f591F12a8B6531Be43E0ccF21cd5fA0E80b0": {
      "balance": "10011000000000000000000"
    },
    "718A8AC0943a6D3FFa3Ec670086bfB03817ed540": {
      "balance": "10011000000000000000000"
    },
    "b30980cE21679314E240DE5Cbf437C15ad459EB8": {
      "balance": "10011000000000000000000"
    },
    "99eCa23623E59C795EceB0edB666eca9eC272339": {
      "balance": "10011000000000000000000"
    },
    "c030e92d19229c3EfD708cf4B85876543ee1A3F7": {
      "balance": "10011000000000000000000"
    },
    "5c98A3414Cb6Ff5c24d145F952Cd19F5f1f56643": {
      "balance": "10011000000000000000000"
    },
    "1979b042Ae2272197f0b74170B3a6F44C3cC5c05": {
      "balance": "10011000000000000000000"
    },
    "Db871070334b961804A15f3606fBB4fAc7C7f932": {
      "balance": "10011000000000000000000"
    },
    "C656C97b765D61E0fbCb1197dC1F3a91CC80C2a4": {
      "balance": "10011000000000000000000"
    },
    "aD95a2f518c197dc9b12eE6381D88bba11F2E0E5": {
      "balance": "10011000000000000000000"
    },
    "4D4B5bF033E4A7359146C9ddb13B1C821FE1D0d3": {
      "balance": "10011000000000000000000"
    },
    "9C64dA169d71C57f85B3d7A17DB27C1ce94FBDE4": {
      "balance": "10011000000000000000000"
    },
    "B5f32e89ccaD3D396f50da32E0a599E43CE87dd7": {
      "balance": "10011000000000000000000"
    },
    "Ba40Db8ab5325494C9E7e07A4c4720990A39305c": {
      "balance": "10011000000000000000000"
    },
    "8B7852DA535df3D06D6ADc1906778afd9481588a": {
      "balance": "10011000000000000000000"
    },
    "a8F41EA062C22dAFFc61e47cF15fc898517b86B1": {
      "balance": "10011000000000000000000"
    },
    "66a3Fc7E8fd6932568cDB6610F5a67BeD9F5beF8": {
      "balance": "10011000000000000000000"
    },
    "10301d9389653497F62876f450332467E07eEe1F": {
      "balance": "10011000000000000000000"
    },
    "6c3ac5fcb13E8DCd908C405Ec6DAcF0EF575D8FC": {
      "balance": "10011000000000000000000"
    },
    "85226637919D3d47E1A37b3AF989E9aE1a1C4790": {
      "balance": "10011000000000000000000"
    },
    "43BCa16603c56cb681d1da3636B7a1A225598bfc": {
      "balance": "10011000000000000000000"
    },
    "E55d8Bc08025BDDF8Da02eEB54882d0586f90700": {
      "balance": "10011000000000000000000"
    },
    "40E1C73f6228a2c15e10aF2F3e890098b777ED15": {
      "balance": "10011000000000000000000"
    },
    "DbbF476089a186a406EA13a4c46813f4BccC3660": {
      "balance": "10011000000000000000000"
    },
    "7baCEA66a75dD974Ad549987768bF8d8908b4917": {
      "balance": "10011000000000000000000"
    },
    "fbF4C2362a9EB672BAC39A46AFd919B3c12Ce44c": {
      "balance": "10011000000000000000000"
    },
    "A8dB96136990be5B3d3bfe592e5A5a5223350A7A": {
      "balance": "10011000000000000000000"
    },
    "1Dd21ED691195EBA816d59B3De7Fab8b3470Ae4B": {
      "balance": "10011000000000000000000"
    },
    "058A778A6aeEfacc013afba92578A43e38cc012D": {
      "balance": "10011000000000000000000"
    },
    "13f52Ab66871880DC8F2179d705281a4cf6a15fB": {
      "balance": "10011000000000000000000"
    },
    "eD1Ed9a71E313d1BCe14aB998E0646F212230a33": {
      "balance": "10011000000000000000000"
    },
    "c563F264f98e34A409C6a085da7510De8B6FE90B": {
      "balance": "10011000000000000000000"
    },
    "c6D678fC6Cc1dA9D5eD1c0075cF7c679e7138e02": {
      "balance": "10011000000000000000000"
    },
    "5179fc80CaB9BB20d5405a50ec0Fb9a36c1B367a": {
      "balance": "10011000000000000000000"
    },
    "0d473f73AAf1C2bf7EBd2be7196C71dBa6C1724b": {
      "balance": "100110000000000000000"
    },
    "6958c5b7E3D94B041d0d76Cac2e09378d31201bd": {
      "balance": "10011000000000000000000"
    },
    "628d4A734d1a2647c67D254209e7B6471a11a5cb": {
      "balance": "10011000000000000000000"
    },
    "E1601e3172F0ef0100e363B639Bd44420B7E5490": {
      "balance": "10011000000000000000000"
    },
    "3337F2Cd103976F044b55D3E69aB06d1ebB142Db": {
      "balance": "10011000000000000000000"
    },
    "8D0D5c57dC232Be15Df4A1a048EF36162C853b94": {
      "balance": "10011000000000000000000"
    },
    "14800c28F3cF1Dd17AaC55263ef4e173b0e8e3Ef": {
      "balance": "10011000000000000000000"
    },
    "f3996A0f0f593BfD5E39780059C5430fab7359FD": {
      "balance": "10011000000000000000000"
    },
    "2217FeBe31Aea6C771AF163dCc453F9f060a4a00": {
      "balance": "10011000000000000000000"
    },
    "f426CC817400766cd6b44F13Cb63Ca648e323484": {
      "balance": "10011000000000000000000"
    },
    "B2C4913e257a34445Ec31685E625bb4060FB8e1f": {
      "balance": "10011000000000000000000"
    },
    "9438dbD05dfC19F049a469185c7599daa82646e8": {
      "balance": "10011000000000000000000"
    },
    "4BeD66Bf507f3CF524704267908Ea4ee3cDe3053": {
      "balance": "10011000000000000000000"
    },
    "9a850fe8105e9CCfBD9d1D06D535BB4948f3f6Cf": {
      "balance": "10011000000000000000000"
    },
    "1277eE554565542A8d0553E1e54006d006db75bd": {
      "balance": "10011000000000000000000"
    },
    "D7e829bE8E374D3fBbd2F68D9A916cB2f769BA89": {
      "balance": "10011000000000000000000"
    },
    "3691b847eD14E296afC90Ff3E37D21e518306170": {
      "balance": "10011000000000000000000"
    },
    "c4C703357B01672cF95bFa0450a5717812Bc7ffb": {
      "balance": "10011000000000000000000"
    },
    "0c9369077836353A8D92aeD29C72A7DfD300B354": {
      "balance": "10011000000000000000000"
    },
    "856DF2A3bdBb8086cE406C469dDE94d12C1E3176": {
      "balance": "10011000000000000000000"
    },
    "E40B3e5c59e2157037b699895329DBe4aA33C039": {
      "balance": "10011000000000000000000"
    },
    "edb47aF3aC2325735722450D1E7DA082bDDad58c": {
      "balance": "10011000000000000000000"
    },
    "315D669866E13fA302B76c85481F9181e06304Ce": {
      "balance": "10011000000000000000000"
    },
    "A5185E3328592428d5989422e0339247dD77e10D": {
      "balance": "10011000000000000000000"
    },
    "85Fd1d1Cd6655EbB89db7D6cA0a5C9c62F7a4CFf": {
      "balance": "10011000000000000000000"
    },
    "ACC9E4430EC1011673547395A191C6b152763EA4": {
      "balance": "10011000000000000000000"
    },
    "3824967C172D52128522dD257FE8f58C9099166B": {
      "balance": "10011000000000000000000"
    },
    "5542aDEA3092da5541250d70a3Db28Ad9BE7Cfc7": {
      "balance": "10011000000000000000000"
    },
    "c61Cd4477f0A98BfC97744481181730f7af7c14f": {
      "balance": "10011000000000000000000"
    },
    "5D7Ffd0fC6DAA67AbF7d48ae69f09dbe53d86983": {
      "balance": "10011000000000000000000"
    },
    "350914ABD4F095534823C1e8fA1cfD7EF79e7E4c": {
      "balance": "10011000000000000000000"
    },
    "ECa6f058B718E320c1D45f5D1fb07947367C3D4B": {
      "balance": "10011000000000000000000"
    },
    "9C577D0795Ed0cA88814d149c2DC61E8Fc48Ad81": {
      "balance": "10011000000000000000000"
    },
    "72fE8bC8E3Ff1e56543c9c1F9834D6dfC31BEDDC": {
      "balance": "10011000000000000000000"
    },
    "6Ff2CFa7899073CD029267fd821C9497811b5f7E": {
      "balance": "10011000000000000000000"
    },
    "4685D123aE928a7912646681ba32035ad6F010a6": {
      "balance": "10011000000000000000000"
    },
    "4799946c8B21fF5E58A225AeCB6F54ec17a94566": {
      "balance": "10011000000000000000000"
    },
    "1D7dA5a23a99Fc33e2e94d502E4Fdb564eA0B24C": {
      "balance": "10011000000000000000000"
    },
    "DFc9719cD9c7982e4A1FFB4B87cC3b861C40E367": {
      "balance": "10011000000000000000000"
    },
    "0c1F0457ce3e87f5eA8F2C3A007dfe963A6Ff9a7": {
      "balance": "10011000000000000000000"
    },
    "7dC23b30dFDc326B9a694c6f9723DC889fe16b7d": {
      "balance": "10011000000000000000000"
    },
    "3F0c4cFDD40D16B7C15878AcCdc91Be9ca4DeE79": {
      "balance": "10011000000000000000000"
    },
    "B984a83416F560437C7866e26CdDb94bDB821594": {
      "balance": "10011000000000000000000"
    },
    "138EA4C57F5b3984EFacd944b3b85dfDd5A78Dcc": {
      "balance": "10011000000000000000000"
    },
    "AD4f16F3435E849505C643714C9E5f40f73c4a5a": {
      "balance": "10011000000000000000000"
    },
    "6b38E861ec0b65fd288d96d5630711C576362152": {
      "balance": "10011000000000000000000"
    },
    "AE15D05100CE807d0aC93119f4ada8fa21441Fd2": {
      "balance": "10011000000000000000000"
    },
    "e0e25c5734bef8b2Add633eAa2518B207DAa0D66": {
      "balance": "10011000000000000000000"
    },
    "9039Ce107A9cD36Ed116958E50f8BDe090e2406f": {
      "balance": "10011000000000000000000"
    },
    "089bE2dD42096ebA1d94aad20228b75df2BeeBC7": {
      "balance": "10011000000000000000000"
    },
    "E3a79AEee437532313015892B52b65f52794F8a2": {
      "balance": "10011000000000000000000"
    },
    "Cc38EE244819649C9DaB02e268306cED09B20672": {
      "balance": "10011000000000000000000"
    },
    "eb0357140a1a0A6c1cB9c93Bf9354ef7365C97d9": {
      "balance": "10011000000000000000000"
    },
    "44370D6b2d010C9eBFa280b6C00010AC99a45660": {
      "balance": "10011000000000000000000"
    },
    "762438915209d038340C3Af9f8aAb8F93aDc8A9A": {
      "balance": "10011000000000000000000"
    },
    "9CBa7aD50fa366Ff6fC2CAe468929eC9AD23Ea2B": {
      "balance": "10011000000000000000000"
    },
    "4f4F159826b2B1eE903A811fCd86E450c9954396": {
      "balance": "10011000000000000000000"
    },
    "3C132B8465e2D172BB7bab6654D85E398ee7c8AD": {
      "balance": "10011000000000000000000"
    },
    "0582426C929B7e525c22201Bd4c143E45189C589": {
      "balance": "10011000000000000000000"
    },
    "fb542740B34dDC0ADE383F2907a1e1E175E0BF5a": {
      "balance": "10011000000000000000000"
    },
    "184Ca91AfE8F36bC5772b29cE2A76c90fCef34D0": {
      "balance": "10011000000000000000000"
    },
    "0C6f48B50B166ddcE52CEE051acCAfFB8ecB4976": {
      "balance": "10011000000000000000000"
    },
    "3aD2bE38fA3DFa7969E79B4053868FD1C368eAb2": {
      "balance": "10011000000000000000000"
    },
    "a6A690637b088E9A1A89c44c9dC5e14eD4825053": {
      "balance": "10011000000000000000000"
    },
    "C224B131Ea71e11E7DF38de3774AAAAe7E197BA4": {
      "balance": "10011000000000000000000"
    },
    "d3C18531f0879B9FB8Ed45830C4ce6b54dC57128": {
      "balance": "10011000000000000000000"
    },
    "02a272d17E1308beF21E783A93D1658f84F2D414": {
      "balance": "10011000000000000000000"
    },
    "57A1aC8167d94b899b32C38Ff9D2B2bD0e55C10d": {
      "balance": "10011000000000000000000"
    },
    "F8fc7D740929E5DD4eBA8fd5a6873Be6a4151087": {
      "balance": "10011000000000000000000"
    },
    "B2AfC45838b364240dE17D3143AA6096d3340A91": {
      "balance": "10011000000000000000000"
    },
    "eAf133d1e0Dd325721665B19f67C9b914EE2469F": {
      "balance": "10011000000000000000000"
    },
    "B7660F1B075e56780e7E026ff66995765f5f1f7F": {
      "balance": "10011000000000000000000"
    },
    "F25087E27B7a59003bb08d2cAc7A69E7c15a4be8": {
      "balance": "10011000000000000000000"
    },
    "E65054681206658A845140331459A057C4EB3CA7": {
      "balance": "10011000000000000000000"
    },
    "e7569A0F93E832a6633d133d23503B5175bEa5Db": {
      "balance": "10011000000000000000000"
    },
    "a9f6102BCf5351dFdC7fA0CA4Fa0A711e16605c3": {
      "balance": "10011000000000000000000"
    },
    "1AB9aA0E855DF953CF8d9cC166172799afD12a68": {
      "balance": "10011000000000000000000"
    },
    "6C04aA35c377E65658EC3600Cab5E8FFa95567D9": {
      "balance": "10011000000000000000000"
    },
    "6b82AD37e64c91c628305813B2DA82F18f8e2a2B": {
      "balance": "10011000000000000000000"
    },
    "AD5D1DeD72F0e70a0a5500B26b82B1A2e8A63471": {
      "balance": "10011000000000000000000"
    },
    "72B3589771Ec8e189a5d9Fe7a214e44085e89054": {
      "balance": "10011000000000000000000"
    },
    "74F57dA8be3E9AB4463DD70319A06Fb5E3168211": {
      "balance": "10011000000000000000000"
    },
    "b6f7F57b99DB21027875BEa3b8531d5925c346cE": {
      "balance": "10011000000000000000000"
    },
    "279d05241d33Dc422d5AEcAc0e089B7f50f879c3": {
      "balance": "10011000000000000000000"
    },
    "d57FEfe1B634ab451a6815Cd6769182EABA62779": {
      "balance": "10011000000000000000000"
    },
    "e86C8538Bdfb253E8D6cC29ee24A330905324849": {
      "balance": "10011000000000000000000"
    },
    "2C58D7f7f9CDF79CF3Cd5F4247761b93428A4E9e": {
      "balance": "10011000000000000000000"
    },
    "37326cEfAFB1676f7Af1CcDcCD37A846Ec64F19d": {
      "balance": "10011000000000000000000"
    },
    "f01DCf91d5f74BDB161F520e800c64F686Eb253F": {
      "balance": "10011000000000000000000"
    },
    "Ba85246bc2A4fdaC1cB2e3C68383Fe79A6466fd9": {
      "balance": "10011000000000000000000"
    },
    "4A76f81eA26381981a3B740975fb4F605989b585": {
      "balance": "10011000000000000000000"
    },
    "00ee7168618BaE4F4d2900D5063c62948c6F0566": {
      "balance": "10011000000000000000000"
    },
    "E1aD0B232B4262E4A279C91070417DAAF202623F": {
      "balance": "10011000000000000000000"
    },
    "f611173319b22080E0F02eE724781d85f4b39Ae6": {
      "balance": "10011000000000000000000"
    },
    "158659458dff3a9E5182cA0e8Ba08F53463FA5e7": {
      "balance": "10011000000000000000000"
    },
    "FEB11610ad367b0c994274A8153E50F4557e473F": {
      "balance": "10011000000000000000000"
    },
    "e1eB2279f45760Ab9D734782B1a0A8FD3d47D807": {
      "balance": "10011000000000000000000"
    },
    "8667d005eCF50Eb247890a11FCdCfC321DC1Da9f": {
      "balance": "10011000000000000000000"
    },
    "5Ce612A664C2f35558Dcab7edb999619e155CD07": {
      "balance": "10011000000000000000000"
    },
    "aD95f88cCd3aBC12ddd6cD0b9a777B95339b747b": {
      "balance": "10011000000000000000000"
    },
    "6E5a5A2963F6d0C2EA26682a152fE3ac7CBC1227": {
      "balance": "10011000000000000000000"
    },
    "000000000000000000000000000000000000ce10": {
      "code": "0x60806040526004361061004a5760003560e01c806303386ba3146101e757806342404e0714610280578063bb913f41146102d7578063d29d44ee14610328578063f7e6af8014610379575b6000600160405180807f656970313936372e70726f78792e696d706c656d656e746174696f6e00000000815250601c019050604051809103902060001c0360001b9050600081549050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610136576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f4e6f20496d706c656d656e746174696f6e20736574000000000000000000000081525060200191505060405180910390fd5b61013f816103d0565b6101b1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f496e76616c696420636f6e74726163742061646472657373000000000000000081525060200191505060405180910390fd5b60405136810160405236600082376000803683855af43d604051818101604052816000823e82600081146101e3578282f35b8282fd5b61027e600480360360408110156101fd57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019064010000000081111561023a57600080fd5b82018360208201111561024c57600080fd5b8035906020019184600183028401116401000000008311171561026e57600080fd5b909192939192939050505061041b565b005b34801561028c57600080fd5b506102956105c1565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156102e357600080fd5b50610326600480360360208110156102fa57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061060d565b005b34801561033457600080fd5b506103776004803603602081101561034b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506107bd565b005b34801561038557600080fd5b5061038e610871565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b60008060007fc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a47060001b9050833f915080821415801561041257506000801b8214155b92505050919050565b610423610871565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146104c3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f73656e64657220776173206e6f74206f776e657200000000000000000000000081525060200191505060405180910390fd5b6104cc8361060d565b600060608473ffffffffffffffffffffffffffffffffffffffff168484604051808383808284378083019250505092505050600060405180830381855af49150503d8060008114610539576040519150601f19603f3d011682016040523d82523d6000602084013e61053e565b606091505b508092508193505050816105ba576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f696e697469616c697a6174696f6e2063616c6c6261636b206661696c6564000081525060200191505060405180910390fd5b5050505050565b600080600160405180807f656970313936372e70726f78792e696d706c656d656e746174696f6e00000000815250601c019050604051809103902060001c0360001b9050805491505090565b610615610871565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146106b5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f73656e64657220776173206e6f74206f776e657200000000000000000000000081525060200191505060405180910390fd5b6000600160405180807f656970313936372e70726f78792e696d706c656d656e746174696f6e00000000815250601c019050604051809103902060001c0360001b9050610701826103d0565b610773576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f496e76616c696420636f6e74726163742061646472657373000000000000000081525060200191505060405180910390fd5b8181558173ffffffffffffffffffffffffffffffffffffffff167fab64f92ab780ecbf4f3866f57cee465ff36c89450dcce20237ca7a8d81fb7d1360405160405180910390a25050565b6107c5610871565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610865576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f73656e64657220776173206e6f74206f776e657200000000000000000000000081525060200191505060405180910390fd5b61086e816108bd565b50565b600080600160405180807f656970313936372e70726f78792e61646d696e000000000000000000000000008152506013019050604051809103902060001c0360001b9050805491505090565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610960576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260118152602001807f6f776e65722063616e6e6f74206265203000000000000000000000000000000081525060200191505060405180910390fd5b6000600160405180807f656970313936372e70726f78792e61646d696e000000000000000000000000008152506013019050604051809103902060001c0360001b90508181558173ffffffffffffffffffffffffffffffffffffffff167f50146d0e3c60aa1d17a70635b05494f864e86144a2201275021014fbf08bafe260405160405180910390a2505056fea165627a7a72305820959a50d5df76f90bc1825042f47788ee27f1b4725f7ed5d37c5c05c0732ef44f0029",
      "storage": {
        "0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103": "0x0Cc59Ed03B3e763c02d54D695FFE353055f1502D"
      },
      "balance": "0"
	}
}`

const alfajoresAllocJSON = `{
  "456f41406B32c45D59E539e4BBA3D7898c3584dA": {
	"balance": "103010030000000000000000000"
  },
  "DD1F519F63423045F526b8c83edC0eB4BA6434a4": {
	"balance": "10011000000000000000000"
  },
  "050f34537F5b2a00B9B9C752Cb8500a3fcE3DA7d": {
	"balance": "10011000000000000000000"
  },
  "Cda518F6b5a797C3EC45D37c65b83e0b0748eDca": {
	"balance": "10011000000000000000000"
  },
  "b4e92c94A2712e98c020A81868264bdE52C188Cb": {
	"balance": "10011000000000000000000"
  },
  "Ae1ec841923811219b98ACeB1db297AADE2F46F3": {
	"balance": "10011000000000000000000"
  },
  "621843731fe33418007C06ee48CfD71e0ea828d9": {
	"balance": "10011000000000000000000"
  },
  "2A43f97f8BF959E31F69A894ebD80A88572C8553": {
	"balance": "10011000000000000000000"
  },
  "AD682035bE6Ab6f06e478D2BDab0EAb6477B460E": {
	"balance": "10011000000000000000000"
  },
  "30D060F129817c4DE5fBc1366d53e19f43c8c64f": {
	"balance": "10011000000000000000000"
  },
  "22579CA45eE22E2E16dDF72D955D6cf4c767B0eF": {
	"balance": "10011000000000000000000"
  },
  "1173C5A50bf025e8356823a068E396ccF2bE696C": {
	"balance": "10011000000000000000000"
  },
  "40F71B525A96baa8d14Eaa7Bcd19929782659c64": {
	"balance": "10011000000000000000000"
  },
  "b923626C6f1d237252793FB2aA12BA21328C51BC": {
	"balance": "10011000000000000000000"
  },
  "B70f9ABf41F36B3ab60cc9aE1a85Ddda3C88D261": {
	"balance": "10011000000000000000000"
  },
  "d4369DB59eaDc4Cfa089c0a3c1004ceAb1b318D8": {
	"balance": "10011000000000000000000"
  },
  "2fd430d3a96eadc38cc1B38b6685C5f52Cf7a083": {
	"balance": "10011000000000000000000"
  },
  "Fecc71C8f33Ca5952534fd346ADdeDC38DBb9cb7": {
	"balance": "10011000000000000000000"
  },
  "0de78C89e7BF5060f28dd3f820C15C4A6A81AFB5": {
	"balance": "10011000000000000000000"
  },
  "75411b92fcE120C1e7fd171b1c2bF802f2E3CF48": {
	"balance": "10011000000000000000000"
  },
  "563433bD8357b06982Fe001df20B2b43393d21d2": {
	"balance": "10011000000000000000000"
  },
  "79dfB9d2367E7921d4139D7841d24ED82F48907F": {
	"balance": "10011000000000000000000"
  },
  "5809369FC5121a071eE67659a975e88ae40fBE3b": {
	"balance": "10011000000000000000000"
  },
  "7517E54a456bcc6c5c695B5d9f97EBc05d29a824": {
	"balance": "10011000000000000000000"
  },
  "B0a1A5Ffcb34E6Fa278D2b40613f0AE1042d32f8": {
	"balance": "10011000000000000000000"
  },
  "EeE9f4DDf49976251E84182AbfD3300Ee58D12aa": {
	"balance": "10011000000000000000000"
  },
  "Eb5Fd57f87a4e1c7bAa53ec1c0d021bb1710B743": {
	"balance": "10011000000000000000000"
  },
  "B7Dd51bFb73c5753778e5Af56f1D9669BCe6777F": {
	"balance": "10011000000000000000000"
  },
  "33C222BB13C63295AF32D6C91278AA34b573e776": {
	"balance": "10011000000000000000000"
  },
  "83c58603bF72DA067D7f6238E7bF390d91B2f531": {
	"balance": "10011000000000000000000"
  },
  "6651112198C0da05921355642a2B8dF1fA3Ede93": {
	"balance": "10011000000000000000000"
  },
  "4EE72A98549eA7CF774C3E2E1b39fF166b4b68BE": {
	"balance": "10011000000000000000000"
  },
  "840b32F30e1a3b2E8b9E6C0972eBa0148E22B847": {
	"balance": "100000000000000000000"
  },
  "000000000000000000000000000000000000ce10": {
	"code": "0x60806040526004361061004a5760003560e01c806303386ba3146101e757806342404e0714610280578063bb913f41146102d7578063d29d44ee14610328578063f7e6af8014610379575b6000600160405180807f656970313936372e70726f78792e696d706c656d656e746174696f6e00000000815250601c019050604051809103902060001c0360001b9050600081549050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610136576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f4e6f20496d706c656d656e746174696f6e20736574000000000000000000000081525060200191505060405180910390fd5b61013f816103d0565b6101b1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f496e76616c696420636f6e74726163742061646472657373000000000000000081525060200191505060405180910390fd5b60405136810160405236600082376000803683855af43d604051818101604052816000823e82600081146101e3578282f35b8282fd5b61027e600480360360408110156101fd57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019064010000000081111561023a57600080fd5b82018360208201111561024c57600080fd5b8035906020019184600183028401116401000000008311171561026e57600080fd5b909192939192939050505061041b565b005b34801561028c57600080fd5b506102956105c1565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156102e357600080fd5b50610326600480360360208110156102fa57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061060d565b005b34801561033457600080fd5b506103776004803603602081101561034b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506107bd565b005b34801561038557600080fd5b5061038e610871565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b60008060007fc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a47060001b9050833f915080821415801561041257506000801b8214155b92505050919050565b610423610871565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146104c3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f73656e64657220776173206e6f74206f776e657200000000000000000000000081525060200191505060405180910390fd5b6104cc8361060d565b600060608473ffffffffffffffffffffffffffffffffffffffff168484604051808383808284378083019250505092505050600060405180830381855af49150503d8060008114610539576040519150601f19603f3d011682016040523d82523d6000602084013e61053e565b606091505b508092508193505050816105ba576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f696e697469616c697a6174696f6e2063616c6c6261636b206661696c6564000081525060200191505060405180910390fd5b5050505050565b600080600160405180807f656970313936372e70726f78792e696d706c656d656e746174696f6e00000000815250601c019050604051809103902060001c0360001b9050805491505090565b610615610871565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146106b5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f73656e64657220776173206e6f74206f776e657200000000000000000000000081525060200191505060405180910390fd5b6000600160405180807f656970313936372e70726f78792e696d706c656d656e746174696f6e00000000815250601c019050604051809103902060001c0360001b9050610701826103d0565b610773576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f496e76616c696420636f6e74726163742061646472657373000000000000000081525060200191505060405180910390fd5b8181558173ffffffffffffffffffffffffffffffffffffffff167fab64f92ab780ecbf4f3866f57cee465ff36c89450dcce20237ca7a8d81fb7d1360405160405180910390a25050565b6107c5610871565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610865576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f73656e64657220776173206e6f74206f776e657200000000000000000000000081525060200191505060405180910390fd5b61086e816108bd565b50565b600080600160405180807f656970313936372e70726f78792e61646d696e000000000000000000000000008152506013019050604051809103902060001c0360001b9050805491505090565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610960576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260118152602001807f6f776e65722063616e6e6f74206265203000000000000000000000000000000081525060200191505060405180910390fd5b6000600160405180807f656970313936372e70726f78792e61646d696e000000000000000000000000008152506013019050604051809103902060001c0360001b90508181558173ffffffffffffffffffffffffffffffffffffffff167f50146d0e3c60aa1d17a70635b05494f864e86144a2201275021014fbf08bafe260405160405180910390a2505056fea165627a7a723058202dbb6037e4381b4ad95015ed99441a23345cc2ae52ef27e2e91d34fb0acd277b0029",
	"storage": {
	  "0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103": "456f41406B32c45D59E539e4BBA3D7898c3584dA"
	},
	"balance": "0"
  }
}`

const mainnetAllocJSON = "{\"0x11901cf7eEae1E2644995FB2E47Ce46bC7F33246\":{\"balance\":\"120000000000000000000000000\"},\"0xC1cDA18694F5B86cFB80c1B4f8Cc046B0d7E6326\":{\"balance\":\"20000000000000000000000000\"},\"0xa5d40D93b01AfBafec84E20018Aff427628F645E\":{\"balance\":\"20000000000000000000000000\"},\"0x8d485780E84E23437f8F6938D96B964645529127\":{\"balance\":\"20000000000000000000000000\"},\"0x5F857c501b73ddFA804234f1f1418D6f75554076\":{\"balance\":\"20000000000000000000000000\"},\"0xaa9064F57F8d7de4b3e08c35561E21Afd6341390\":{\"balance\":\"20000000000000000000000000\"},\"0x7FA26b50b3e9a2eC8AD1850a4c4FBBF94D806E95\":{\"balance\":\"20000000000000000000000000\"},\"0x08960Ce6b58BE32FBc6aC1489d04364B4f7dC216\":{\"balance\":\"20000000000000000000000000\"},\"0x77B68B2e7091D4F242a8Af89F200Af941433C6d8\":{\"balance\":\"20000000000000000000000000\"},\"0x75Bb69C002C43f5a26a2A620518775795Fd45ecf\":{\"balance\":\"20000000000000000000000000\"},\"0x19992AE48914a178Bf138665CffDD8CD79b99513\":{\"balance\":\"20000000000000000000000000\"},\"0xE23a4c6615669526Ab58E9c37088bee4eD2b2dEE\":{\"balance\":\"20000000000000000000000\"},\"0xDe22679dCA843B424FD0BBd70A22D5F5a4B94fe4\":{\"balance\":\"10200014000000000000000000\"},\"0x743D80810fe10c5C3346D2940997cC9647035B13\":{\"balance\":\"20513322000000000000000000\"},\"0x8e1c4355307F1A59E7eD4Ae057c51368b9338C38\":{\"balance\":\"7291740000000000000000000\"},\"0x417fe63186C388812e342c85FF87187Dc584C630\":{\"balance\":\"20000062000000000000000000\"},\"0xF5720c180a6Fa14ECcE82FB1bB060A39E93A263c\":{\"balance\":\"30000061000000000000000000\"},\"0xB80d1e7F9CEbe4b5E1B1Acf037d3a44871105041\":{\"balance\":\"9581366833333333333333335\"},\"0xf8ed78A113cD2a34dF451Ba3D540FFAE66829AA0\":{\"balance\":\"11218686833333333333333333\"},\"0x9033ff75af27222c8f36a148800c7331581933F3\":{\"balance\":\"11218686833333333333333333\"},\"0x8A07541C2eF161F4e3f8de7c7894718dA26626B2\":{\"balance\":\"11218686833333333333333333\"},\"0xB2fe7AFe178335CEc3564d7671EEbD7634C626B0\":{\"balance\":\"11218686833333333333333333\"},\"0xc471776eA02705004C451959129bF09423B56526\":{\"balance\":\"11218686833333333333333333\"},\"0xeF283eca68DE87E051D427b4be152A7403110647\":{\"balance\":\"14375000000000000000000000\"},\"0x7cf091C954ed7E9304452d31fd59999505Ddcb7a\":{\"balance\":\"14375000000000000000000000\"},\"0xa5d2944C32a8D7b284fF0b84c20fDcc46937Cf64\":{\"balance\":\"14375000000000000000000000\"},\"0xFC89C17525f08F2Bc9bA8cb77BcF05055B1F7059\":{\"balance\":\"14375000000000000000000000\"},\"0x3Fa7C646599F3174380BD9a7B6efCde90b5d129d\":{\"balance\":\"14375000000000000000000000\"},\"0x989e1a3B344A43911e02cCC609D469fbc15AB1F1\":{\"balance\":\"14375000000000000000000000\"},\"0xAe1d640648009DbE0Aa4485d3BfBB68C37710924\":{\"balance\":\"20025000000000000000000000\"},\"0x1B6C64779F42BA6B54C853Ab70171aCd81b072F7\":{\"balance\":\"20025000000000000000000000\"},\"000000000000000000000000000000000000ce10\":{\"code\":\"0x60806040526004361061004a5760003560e01c806303386ba3146101e757806342404e0714610280578063bb913f41146102d7578063d29d44ee14610328578063f7e6af8014610379575b6000600160405180807f656970313936372e70726f78792e696d706c656d656e746174696f6e00000000815250601c019050604051809103902060001c0360001b9050600081549050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610136576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f4e6f20496d706c656d656e746174696f6e20736574000000000000000000000081525060200191505060405180910390fd5b61013f816103d0565b6101b1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f496e76616c696420636f6e74726163742061646472657373000000000000000081525060200191505060405180910390fd5b60405136810160405236600082376000803683855af43d604051818101604052816000823e82600081146101e3578282f35b8282fd5b61027e600480360360408110156101fd57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019064010000000081111561023a57600080fd5b82018360208201111561024c57600080fd5b8035906020019184600183028401116401000000008311171561026e57600080fd5b909192939192939050505061041b565b005b34801561028c57600080fd5b506102956105c1565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156102e357600080fd5b50610326600480360360208110156102fa57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061060d565b005b34801561033457600080fd5b506103776004803603602081101561034b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506107bd565b005b34801561038557600080fd5b5061038e610871565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b60008060007fc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a47060001b9050833f915080821415801561041257506000801b8214155b92505050919050565b610423610871565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146104c3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f73656e64657220776173206e6f74206f776e657200000000000000000000000081525060200191505060405180910390fd5b6104cc8361060d565b600060608473ffffffffffffffffffffffffffffffffffffffff168484604051808383808284378083019250505092505050600060405180830381855af49150503d8060008114610539576040519150601f19603f3d011682016040523d82523d6000602084013e61053e565b606091505b508092508193505050816105ba576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f696e697469616c697a6174696f6e2063616c6c6261636b206661696c6564000081525060200191505060405180910390fd5b5050505050565b600080600160405180807f656970313936372e70726f78792e696d706c656d656e746174696f6e00000000815250601c019050604051809103902060001c0360001b9050805491505090565b610615610871565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146106b5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f73656e64657220776173206e6f74206f776e657200000000000000000000000081525060200191505060405180910390fd5b6000600160405180807f656970313936372e70726f78792e696d706c656d656e746174696f6e00000000815250601c019050604051809103902060001c0360001b9050610701826103d0565b610773576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260188152602001807f496e76616c696420636f6e74726163742061646472657373000000000000000081525060200191505060405180910390fd5b8181558173ffffffffffffffffffffffffffffffffffffffff167fab64f92ab780ecbf4f3866f57cee465ff36c89450dcce20237ca7a8d81fb7d1360405160405180910390a25050565b6107c5610871565b73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610865576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f73656e64657220776173206e6f74206f776e657200000000000000000000000081525060200191505060405180910390fd5b61086e816108bd565b50565b600080600160405180807f656970313936372e70726f78792e61646d696e000000000000000000000000008152506013019050604051809103902060001c0360001b9050805491505090565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610960576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260118152602001807f6f776e65722063616e6e6f74206265203000000000000000000000000000000081525060200191505060405180910390fd5b6000600160405180807f656970313936372e70726f78792e61646d696e000000000000000000000000008152506013019050604051809103902060001c0360001b90508181558173ffffffffffffffffffffffffffffffffffffffff167f50146d0e3c60aa1d17a70635b05494f864e86144a2201275021014fbf08bafe260405160405180910390a2505056fea165627a7a723058206808dd43e7d765afca53fe439122bc5eac16d708ce7d463451be5042426f101f0029\",\"storage\":{\"0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103\":\"0xE23a4c6615669526Ab58E9c37088bee4eD2b2dEE\"},\"balance\":\"0\"}}"

var celoL1GenesisAllocJSON = map[uint64]string{
	MainnetNetworkID:   mainnetAllocJSON,
	AlfajoresNetworkID: alfajoresAllocJSON,
	BaklavaNetworkID:   baklavaAllocJSON,
}

// GetCeloL1GenesisAlloc returns the legacy Celo L1 genesis allocation JSON for the given network ID.
func GetCeloL1GenesisAlloc(config *params.ChainConfig) ([]byte, error) {
	chainID := config.ChainID.Uint64()
	allocJSON, ok := celoL1GenesisAllocJSON[chainID]
	if !ok {
		return nil, fmt.Errorf("no genesis allocation JSON found for network ID %d", chainID)
	}
	return []byte(allocJSON), nil
}

// BuildGenesis creates a genesis block from the given parameters.
func BuildGenesis(config *params.ChainConfig, allocs, extraData []byte, timestamp uint64) *core.Genesis {
	genesisAlloc := &types.GenesisAlloc{}
	genesisAlloc.UnmarshalJSON(allocs)
	return &core.Genesis{
		Config:    config,
		Timestamp: timestamp,
		ExtraData: extraData,
		Alloc:     *genesisAlloc,
	}
}
