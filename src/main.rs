#[macro_use] extern crate log;
extern crate env_logger;
#[macro_use] extern crate clap;
mod channel;

use clap::{App, Arg};
use channel::{Settings, DuplexChannel};
use std::thread;
use std::time::Duration;

fn main() {
    env_logger::init().unwrap();
    let matches = App::new("riftd")
        .version(crate_version!())
        .author(crate_authors!())
        .about("The Rift protocol daemon")
        .arg(Arg::with_name("DEVICE_NAME")
             .short("s")
             .long("device")
             .takes_value(true)
             .required(true)
             .help("Sets the serial device to use"))
        .arg(Arg::with_name("BROADCAST")
             .short("b")
             .long("broadcast")
             .takes_value(true)
             .default_value("1000")
             .help("Sets the node broadcast interval"))
        .arg(Arg::with_name("POLL_INTERVAL")
             .short("p")
             .long("poll")
             .takes_value(true)
             .default_value("100")
             .help("Sets the I/O event poll interval"))
        .get_matches();

    let device_name = matches.value_of("DEVICE_NAME").unwrap().to_string();
    let broadcast_interval = Duration::from_millis(matches.value_of("BROADCAST").unwrap().parse().unwrap());
    let poll_interval = Duration::from_millis(matches.value_of("POLL_INTERVAL").unwrap().parse().unwrap());

    let (mut channel, tx, rx) = DuplexChannel::open(Settings {
        device_name: device_name
    });
    let reactor = thread::spawn(move || {
        channel.run_loop();
    });
    let receiver = thread::spawn(move || {
        loop {
            match rx.try_recv() {
                Ok(data) => {
                    info!("Received: {:?}", data);
                },
                _ => {}
            }
            thread::sleep(poll_interval);
        }
    });
    let sender = thread::spawn(move || {
        loop {
            let data = vec![126, 0, 16, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 4, 72, 69, 76, 76, 79, 137];
            tx.send(data).unwrap();
            thread::sleep(broadcast_interval);
        }
    });
    reactor.join().unwrap();
    receiver.join().unwrap();
    sender.join().unwrap();
}
