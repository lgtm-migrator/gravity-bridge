//! Basic utility functions to stubbornly get data
use cosmos_gravity::query::get_last_event_nonce;
use ethereum_gravity::types::EthClient;
use ethers::prelude::*;
use ocular::{cosmrs::AccountId, GrpcClient};
use std::time::Duration;
use tokio::time::sleep as delay_for;

pub const RETRY_TIME: Duration = Duration::from_secs(5);

/// gets the current block number, no matter how long it takes
pub async fn get_block_number_with_retry(eth_client: EthClient) -> U64 {
    let mut res = eth_client.get_block_number().await;
    while res.is_err() {
        error!("Failed to get latest block! Is your Eth node working?");
        delay_for(RETRY_TIME).await;
        res = eth_client.get_block_number().await;
    }
    res.unwrap()
}

/// gets the last event nonce, no matter how long it takes.
pub async fn get_last_event_nonce_with_retry(
    cosmos_client: &mut GrpcClient,
    our_cosmos_address: &AccountId,
) -> u64 {
    let mut res = get_last_event_nonce(cosmos_client, our_cosmos_address).await;
    while res.is_err() {
        error!(
            "Failed to get last event nonce, is the Cosmos GRPC working? {:?}",
            res
        );
        delay_for(RETRY_TIME).await;
        res = get_last_event_nonce(cosmos_client, our_cosmos_address).await;
    }
    res.unwrap()
}

/// gets the chain ID, no matter how long it takes
pub async fn get_chain_id_with_retry(eth_client: EthClient) -> U256 {
    let mut res = eth_client.get_chainid().await;
    while res.is_err() {
        error!("Failed to get chain ID! Is your Eth node working?");
        delay_for(RETRY_TIME).await;
        res = eth_client.get_chainid().await;
    }
    res.unwrap()
}
