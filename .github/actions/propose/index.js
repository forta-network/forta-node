const { AdminClient } = require('defender-admin-client');
const core = require('@actions/core');

function createReleaseProposalParams(versionContract, network, multisig, version, cid) {
    return {
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
}

function createPrereleaseProposalParams(versionContract, network, multisig, version, cid) {
    return {
        contract: { address: versionContract, network: network },
        title: `Forta Node Prerelease ${version}`,
        description: `Prerelease forta-node ${version} (${cid})`,
        type: 'custom',
        functionInterface: {
            "inputs": [
                {
                    "internalType": "string",
                    "name": "version",
                    "type": "string"
                }
            ],
            "name": "setScannerNodeBetaVersion",
            "outputs": [],
            "stateMutability": "nonpayable",
            "type": "function"
        },
        functionInputs: [cid],
        via: `${multisig}`,
        viaType: 'Gnosis Safe',
    }    
}

async function createProposal(apiKey, apiSecret, params) {
    const client = new AdminClient({apiKey, apiSecret});
    const result = await client.createProposal(params);
    return result.url
}

async function main(){
    try {
        const apiKey = core.getInput('api-key');
        const apiSecret = core.getInput('api-secret');
        const versionContract = core.getInput('scanner-version-contract');
        const network = core.getInput('network');
        const multisig = core.getInput('multisig');
        const version = core.getInput('version');
        const releaseCid = core.getInput('release-cid');
        const isRelease = core.getInput('is-release');

        const releaseProposalParams = createReleaseProposalParams(versionContract, network, multisig, version, releaseCid)
        const prereleaseProposalParams = createPrereleaseProposalParams(versionContract, network, multisig, version, releaseCid)

        // if we are making a prerelease, we should write the new ref.
        // if we are making a release, we should make sure that the prerelease ref is up-to-date.
        const prereleaseProposalUrl = await createProposal(apiKey, apiSecret, prereleaseProposalParams)
        console.log(`prerelease proposal created: ${prereleaseProposalUrl}`);
        core.setOutput("prerelease-proposal-url", prereleaseProposalUrl);

        // write the release value only if we are making a release
        let releaseProposalUrl = '';
        if (isRelease == 'true') {
            releaseProposalUrl = await createProposal(apiKey, apiSecret, releaseProposalParams)
            console.log(`release proposal created: ${releaseProposalUrl}`);
        }
        core.setOutput("release-proposal-url", releaseProposalUrl);
    } catch (error) {
        core.setFailed(error.message);
    }
}

main().then((url) => {
    console.log(url)
}).catch((e)=>{
    console.log(e)
})
