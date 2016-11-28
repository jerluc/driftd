use mio::*;
use mio::unix::EventedFd;
use mio::channel::{channel, Sender, Receiver};
use mio_uds::UnixDatagram;
use uuid::Uuid;
use std::io::prelude::*;
use std::os::unix::io::AsRawFd;
use std::io::Result;
use std::mem::transmute;
use std::path::{Path, PathBuf};
use std::thread::{JoinHandle, spawn};

const READ: Token = Token(0);
const SERVER: Token = Token(1);

#[derive(Debug)]
pub enum ClientEvent {
    Nop,
    Open(u64),
    Close(u64)
}

pub struct SocketManager {
    server_sock: UnixDatagram,
    sock_dir: PathBuf,
    client_threads: Vec<JoinHandle<()>>
}

impl SocketManager {
    pub fn new<P: AsRef<Path>>(sock_dir_path: P) -> Result<Self> {
        let server_sock_path = sock_dir_path.as_ref().join("riftd");
        let server_sock = try!(UnixDatagram::bind(server_sock_path));
        let sm = SocketManager {
            server_sock: server_sock,
            sock_dir: sock_dir_path.as_ref().to_path_buf(),
            client_threads: Vec::new(),
        };
        Ok(sm)
    }

    fn on_new_client(&mut self, dest_addr: Vec<u8>) -> Result<i32> {
        let unique_id = Uuid::new_v4().to_string();
        let client_sock_path = self.sock_dir.join("riftd-".to_string() + unique_id.as_ref());
        println!("{:?}", client_sock_path);
        let client_sock = try!(UnixDatagram::bind(client_sock_path));
        println!("Created socket");
        let client_fd = client_sock.as_raw_fd();

        let poll = Poll::new().unwrap();
        let mut events = Events::with_capacity(1024);

        try!(poll.register(&client_sock, READ, Ready::readable(), PollOpt::level()));

        let thread = spawn(move || {
            loop {
                match poll.poll(&mut events, None) {
                    Ok(num_events) if num_events > 0 => {
                        for event in events.iter() {
                            match event.token() {
                                READ => {
                                    println!("Got client read event!");
                                    let mut buf = [0; 128];
                                    client_sock.recv(&mut buf).unwrap();
                                    println!("Received: {:?}", buf.to_vec());
                                },
                                _ => unreachable!()
                            }
                        }
                    },
                    _ => unreachable!()
                }
            }
        });
        self.client_threads.push(thread);
        Ok(client_fd)
    }

    pub fn start(&mut self) -> Result<()> {
        let poll = try!(Poll::new());
        let mut events = Events::with_capacity(1024);

        try!(poll.register(&self.server_sock, SERVER, Ready::all(), PollOpt::level()));
        
        loop {
            match poll.poll(&mut events, None) {
                Ok(num_events) if num_events > 0 => {
                    for event in events.iter() {
                        match event.token() {
                            SERVER => {

                                if event.kind().is_readable() { 
                                    let mut buf = [0; 8];
                                    let (_, addr) = try!(self.server_sock.recv_from(&mut buf));
                                    println!("Received client socket request from ({:?}): {:?}", addr, buf);
                                    let client_fd = try!(self.on_new_client(buf.to_vec()));
                                    let client_addr = addr.as_pathname().unwrap();
                                    let client_fd_bytes = fd_to_bytes(client_fd);
                                    println!("Sending back client fd: {:?}", client_fd);
                                    try!(self.server_sock.send_to(&client_fd_bytes, client_addr));
                                } else if !event.kind().is_writable() {
                                    println!("Got event: {:?}", event.kind());
                                }
                            },
                            _ => unreachable!()
                        }
                    }
                },
                _ => unreachable!()
            }
        }

        Ok(())
    }
}

fn fd_to_bytes(fd: i32) -> [u8; 8] {
    let raw_bytes : [u8; 8] = unsafe { transmute(fd as u64) };
    raw_bytes
}
