use crate::application::APP;
use abscissa_core::{clap::Parser, Application, Command, Runnable};

/// Show a Cosmos Key
#[derive(Command, Debug, Default, Parser)]
pub struct ShowCosmosKeyCmd {
    pub args: Vec<String>,
}

// Entry point for `gorc keys cosmos show [name]`
impl Runnable for ShowCosmosKeyCmd {
    fn run(&self) {
        let config = APP.config();
        let name = self.args.get(0).expect("name is required");
        let key = config.load_account(name.clone());

        let address = key
            .address(config.cosmos.prefix.trim())
            .expect("Could not generate public key");

        println!("{}\t{}", name, address)
    }
}
