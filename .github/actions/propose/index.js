const { AdminClient } = require('defender-admin-client');
const core = require('@actions/core');

async function proposeUpgrade(apiKey, apiSecret, versionContract, network, multisig, version, cid) {
    const client = new AdminClient({apiKey, apiSecret});

    console.log({apiKey, apiSecret, versionContract, network, multisig, version, cid})

    const params = {
        contract: { address: versionContract, network: network },
        title: `Forta Node Release ${version}`,
        description: `Release forta-node ${version} (${cid})`,
        type: 'custom',
        functionInterface: {
            "inputs": [
                {
                    "internalType": "string",
                    "name": "version",
                    "type": "string"
                }
            ],
            "name": "setScannerNodeVersion",
            "outputs": [],
            "stateMutability": "nonpayable",
            "type": "function"
        },
        functionInputs: [cid],
        via: `${multisig}`,
        viaType: 'Gnosis Safe',
    }

    const result = await client.createProposal(params);
    return result.url
}

async function main(){
    try {
        const proposalUrl = await proposeUpgrade(
            core.getInput('api-key'),
            core.getInput('api-secret'),
            core.getInput('scanner-version-contract'),
            core.getInput('network'),
            core.getInput('multisig'),
            core.getInput('version'),
            core.getInput('release-cid'))

        console.log(`proposal created: ${proposalUrl}`);
        core.setOutput("proposal-url", proposalUrl);
    } catch (error) {
        core.setFailed(error.message);
    }
}

main().then((url) => {
    console.log(url)
}).catch((e)=>{
    console.log(e)
})