//! ssh-microvm entrypoint.
//!
//! For now this is a stub: it parses CLI flags and prints the resolved config.
//! Subsequent steps wire in the Firecracker API client, the russh server,
//! and the session proxy.

use anyhow::Result;
use clap::Parser;

mod config;

#[tokio::main]
async fn main() -> Result<()> {
    tracing_subscriber::fmt()
        .with_env_filter(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| tracing_subscriber::EnvFilter::new("info")),
        )
        .init();

    let cfg = config::Config::parse();
    tracing::debug!(host_key = ?cfg.host_key_path(), "resolved host key path");
    tracing::info!(?cfg, "ssh-microvm starting (stub)");
    println!("ssh-microvm: skeleton build; nothing wired yet.");
    Ok(())
}
