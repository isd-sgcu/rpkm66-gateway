#[derive(serde::Deserialize)]
pub struct Config {
    pub app: AppConfig,
    pub service: ServiceConfig,
}

#[derive(serde::Deserialize)]
pub struct AppConfig {
    pub port: i32,
    pub debug: bool,
    pub max_file_size: i32,
    pub phase: String,
}

#[derive(serde::Deserialize)]
pub struct ServiceConfig {
    pub backend: String,
    pub auth: String,
    pub file: String,
}

pub fn get_config() -> Config {
    let config = config::Config::builder()
        .add_source(config::File::with_name("config/config"))
        .build()
        .expect("Unable to read config file");

    config
        .try_deserialize()
        .expect("Unable to parse config file")
}
