#!/usr/bin/env bash

# This script is used to generate the getting-started.json configuration file
# used in the Getting Started quickstart guide on the docs site. Avoids the
# need to have the getting-started.json committed to the repo since it's an
# invalid JSON file when not filled in, which is annoying.

reqenv() {
    if [ -z "${!1}" ]; then
        echo "Error: environment variable '$1' is undefined"
        exit 1
    fi
}

# Check required environment variables
reqenv "DEPLOYMENT_CONTEXT"
reqenv "GS_ADMIN_ADDRESS"
reqenv "GS_BATCHER_ADDRESS"
reqenv "GS_PROPOSER_ADDRESS"
reqenv "GS_SEQUENCER_ADDRESS"
reqenv "L1_RPC_URL"
reqenv "L1_CHAIN_ID"
reqenv "L2_CHAIN_ID"
reqenv "L1_BLOCK_TIME"
reqenv "L2_BLOCK_TIME"
reqenv "FUNDS_DEV_ACCOUNTS"
reqenv "USE_PLASMA"
reqenv "DEPLOY_CELO_CONTRACTS"

# Get the finalized block timestamp and hash
block=$(cast block finalized --rpc-url "$L1_RPC_URL")
timestamp=$(echo "$block" | awk '/timestamp/ { print $2 }')
blockhash=$(echo "$block" | awk '/hash/ { print $2 }')
batchInboxAddressSuffix=$(printf "%0$((38 - ${#L2_CHAIN_ID}))d" 0)$L2_CHAIN_ID
batchInboxAddress=0xff$batchInboxAddressSuffix

# Generate the config file
config=$(cat << EOL
{
  "l1StartingBlockTag": "$blockhash",

  "l1ChainID": $L1_CHAIN_ID,
  "l2ChainID": $L2_CHAIN_ID,
  "l2BlockTime": $L2_BLOCK_TIME,
  "l1BlockTime": $L1_BLOCK_TIME,

  "maxSequencerDrift": 600,
  "sequencerWindowSize": 3600,
  "channelTimeout": 300,

  "p2pSequencerAddress": "$GS_SEQUENCER_ADDRESS",
  "batchInboxAddress": "$batchInboxAddress",
  "batchSenderAddress": "$GS_BATCHER_ADDRESS",

  "l2OutputOracleSubmissionInterval": 120,
  "l2OutputOracleStartingBlockNumber": 0,
  "l2OutputOracleStartingTimestamp": $timestamp,

  "l2OutputOracleProposer": "$GS_PROPOSER_ADDRESS",
  "l2OutputOracleChallenger": "$GS_ADMIN_ADDRESS",

  "finalizationPeriodSeconds": 12,

  "proxyAdminOwner": "$GS_ADMIN_ADDRESS",
  "baseFeeVaultRecipient": "$GS_ADMIN_ADDRESS",
  "l1FeeVaultRecipient": "$GS_ADMIN_ADDRESS",
  "sequencerFeeVaultRecipient": "$GS_ADMIN_ADDRESS",
  "finalSystemOwner": "$GS_ADMIN_ADDRESS",
  "superchainConfigGuardian": "$GS_ADMIN_ADDRESS",

  "baseFeeVaultMinimumWithdrawalAmount": "0x8ac7230489e80000",
  "l1FeeVaultMinimumWithdrawalAmount": "0x8ac7230489e80000",
  "sequencerFeeVaultMinimumWithdrawalAmount": "0x8ac7230489e80000",
  "baseFeeVaultWithdrawalNetwork": 0,
  "l1FeeVaultWithdrawalNetwork": 0,
  "sequencerFeeVaultWithdrawalNetwork": 0,

  "gasPriceOracleOverhead": 0,
  "gasPriceOracleScalar": 1000000,

  "deployCeloContracts": $DEPLOY_CELO_CONTRACTS,

  "enableGovernance": $ENABLE_GOVERNANCE,
  "governanceTokenSymbol": "OP",
  "governanceTokenName": "Optimism",
  "governanceTokenOwner": "$GS_ADMIN_ADDRESS",

  "l2GenesisBlockGasLimit": "0x1c9c380",
  "l2GenesisBlockBaseFeePerGas": "0x3b9aca00",
  "l2GenesisRegolithTimeOffset": "0x0",

  "eip1559Denominator": 50,
  "eip1559DenominatorCanyon": 250,
  "eip1559Elasticity": 6,

  "l2GenesisFjordTimeOffset": "0x0",
  "l2GenesisEcotoneTimeOffset": "0x0",
  "l2GenesisDeltaTimeOffset": "0x0",
  "l2GenesisCanyonTimeOffset": "0x0",

  "systemConfigStartBlock": 0,

  "requiredProtocolVersion": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "recommendedProtocolVersion": "0x0000000000000000000000000000000000000000000000000000000000000000",

  "faultGameAbsolutePrestate": "0x03c7ae758795765c6664a5d39bf63841c71ff191e9189522bad8ebff5d4eca98",
  "faultGameMaxDepth": 44,
  "faultGameClockExtension": 0,
  "faultGameMaxClockDuration": 600,
  "faultGameGenesisBlock": 0,
  "faultGameGenesisOutputRoot": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "faultGameSplitDepth": 14,
  "faultGameWithdrawalDelay": 604800,

  "preimageOracleMinProposalSize": 1800000,
  "preimageOracleChallengePeriod": 86400,

  "fundDevAccounts": $FUNDS_DEV_ACCOUNTS,
  "useFaultProofs": false,
  "proofMaturityDelaySeconds": 604800,
  "disputeGameFinalityDelaySeconds": 302400,
  "respectedGameType": 0,

  "usePlasma": $USE_PLASMA,
  "daCommitmentType": "GenericCommitment",
  "daChallengeWindow": 3600,
  "daResolveWindow": 3600,
  "daBondSize": 1000000,
  "daResolverRefundPercentage": 0
}
EOL
)

# Write the config file
echo "$config" > "deploy-config/$DEPLOYMENT_CONTEXT.json"

echo "Wrote config file to deploy-config/$DEPLOYMENT_CONTEXT.json"
