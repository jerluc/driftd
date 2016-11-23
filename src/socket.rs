use notify::{self, RecommendedWatcher, Watcher, RecursiveMode};
use std::sync::mpsc::{channel, Receiver, RecvError};
use std::time::Duration;

pub trait ClientHandler {
    fn on_client_event(&self, event: ClientEvent);
}

#[derive(Debug)]
pub enum ClientEvent {
    Nop,
    Open(u64),
    Close(u64)
}

pub struct SocketWatcher {
    socket_dir: String,
    watcher: RecommendedWatcher,
    client_events: Receiver<notify::DebouncedEvent>
}

impl SocketWatcher {
    pub fn new(socket_dir: &str) -> notify::Result<SocketWatcher> {
        let (tx, rx) = channel();

        let mut watcher: RecommendedWatcher = try!(Watcher::new(tx, Duration::from_secs(1)));
        
        try!(watcher.watch(socket_dir, RecursiveMode::NonRecursive));

        Ok(SocketWatcher {
            socket_dir: socket_dir.to_string(),
            watcher: watcher,
            client_events: rx
        })
    }

    pub fn recv(&self) -> Result<ClientEvent, RecvError> {
        match self.client_events.recv() {
            Ok(event) => match event {
                notify::DebouncedEvent::Create(path) =>
                    Ok(ClientEvent::Open(1)),
                notify::DebouncedEvent::Remove(path) =>
                    Ok(ClientEvent::Close(1)),
                _ => Ok(ClientEvent::Nop)
            },
            Err(e) => Err(e)
        }
    }
}
