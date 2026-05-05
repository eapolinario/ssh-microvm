//! ssh-microvm entrypoint.
//!
//! The dry-boot path exercises the Firecracker boot sequence used by later
//! lifecycle code; the server path accepts SSH connections and validates the
//! outer SSH stack.

use anyhow::Result;
use clap::Parser;
use ssh_microvm::{
    boot,
    config::{Config, RunMode},
    ssh_server,
};

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt()
        .with_env_filter(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| tracing_subscriber::EnvFilter::new("info")),
        )
        .init();

    let cfg = Config::parse();
    tracing::debug!(host_key = ?cfg.host_key_path(), "resolved host key path");
    tracing::info!(?cfg, "ssh-microvm starting");

    match cfg.run_mode() {
        RunMode::BootOnce => boot::dry_boot(&cfg).await?,
        RunMode::Server => ssh_server::run(&cfg).await?,
    }

    Ok(())
}
