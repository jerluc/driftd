extern crate serial;
extern crate mio;

use self::mio::*;
use self::mio::channel::{channel, Sender, Receiver};
use self::mio::unix::EventedFd;
use std::io::prelude::*;
use std::os::unix::io::AsRawFd;

const READ: Token = Token(0);
const WRITE: Token = Token(1);

pub struct Settings {
    pub device_name: String,
}

pub struct DuplexChannel {
    settings: Settings,
    device: serial::SystemPort,
    tx_from_client: Receiver<Vec<u8>>,
    rx_to_client: Sender<Vec<u8>>
}

impl DuplexChannel {
    pub fn open(settings: Settings) -> (DuplexChannel, Sender<Vec<u8>>, Receiver<Vec<u8>>) {
        let device = serial::open(&settings.device_name).unwrap();
        let (client_tx, tx_from_client): (Sender<Vec<u8>>, Receiver<Vec<u8>>) = channel();
        let (rx_to_client, client_rx): (Sender<Vec<u8>>, Receiver<Vec<u8>>) = channel();
        let channel = DuplexChannel {
            settings: settings,
            device: device,
            tx_from_client: tx_from_client,
            rx_to_client: rx_to_client
        };
        (channel, client_tx, client_rx)
    }

    pub fn run_loop(&mut self) {
        trace!("Creating Poll and events queue");
        let poll = Poll::new().unwrap();
        let mut events = Events::with_capacity(1024);

        trace!("Registering IO events");
        let fd = self.device.as_raw_fd();
        let io = EventedFd(&fd);

        poll.register(&io, READ, Ready::readable(), PollOpt::edge()).unwrap();
        poll.register(&self.tx_from_client, WRITE, Ready::readable(), PollOpt::edge()).unwrap();

        // Main reactor loop
        loop {
            // TODO: What good does this timeout value really do?
            match poll.poll(&mut events, None) {
                Ok(num_events) if num_events > 0 => {
                    for event in events.iter() {
                        trace!("Got token: {:?}", event.token());
                        match event.token() {
                            READ => {
                                let mut buf: Vec<u8> = Vec::new();
                                self.device.read_to_end(&mut buf);
                                if buf.len() > 0 {
                                    self.rx_to_client.send(buf);
                                }
                            },
                            WRITE => {
                                match self.tx_from_client.try_recv() {
                                    Ok(outgoing) => {
                                        match self.device.write(&outgoing) {
                                            Ok(num_bytes) => {
                                                info!("Wrote {} bytes", num_bytes);
                                            },
                                            _ => {
                                                error!("Failed to send data");
                                                return;
                                            }
                                        }
                                    },
                                    _ => {}
                                }
                            },
                            _ => {
                                error!("WHAT");
                                return;
                            }
                        }
                    }
                },
                Ok(_) => {},
                Err(e) => {
                    error!("Failed to poll events: {:?}", e);
                    break;
                }
            }
        }
    }
}
