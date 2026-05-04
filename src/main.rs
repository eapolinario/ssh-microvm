//! ssh-microvm entrypoint.
//!
//! For now the server path is a stub. The dry-boot path exercises the
//! Firecracker boot sequence used by later lifecycle code.

use anyhow::Result;
use clap::Parser;
use ssh_microvm::{
    boot,
    config::{Config, RunMode},
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
        RunMode::Server => println!("ssh-microvm: skeleton server; nothing wired yet."),
    }

    Ok(())
}
