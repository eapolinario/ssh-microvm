//! Placeholder for `Config`. Real implementation lands in step 2.
//!
//! Kept minimal so `cargo check` succeeds while we scaffold.

use clap::Parser;

#[derive(Debug, Parser)]
#[command(name = "ssh-microvm", version, about)]
pub struct Config {}
