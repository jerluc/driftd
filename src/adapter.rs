// TODO: Figure out how we can make this useful

use std::sync::mpsc::{self};
use std::marker::Send;
use std::result::Result;
use std::sync::mpsc::{SendError, TryRecvError};

use mio::channel::{self as mioc};

pub trait ISender<T: Send + Copy> {
    fn isend(&self, t: T) -> Result<(), SendError<T>>;
}

pub trait IReceiver<T: Send> {
    fn irecv(&self) -> Result<T, TryRecvError>;
}

impl<T: Send> IReceiver<T> for mpsc::Receiver<T> {
    fn irecv(&self) -> Result<T, TryRecvError> {
        self.try_recv()
    }
}

impl<T: Send + Copy> ISender<T> for mpsc::Sender<T> {
    fn isend(&self, t: T) -> Result<(), SendError<T>> {
        self.send(t)
    }
}

impl<T: Send> IReceiver<T> for mioc::Receiver<T> {
    fn irecv(&self) -> Result<T, TryRecvError> {
        self.try_recv()
    }
}

impl<T: Send + Copy> ISender<T> for mioc::Sender<T> {
    fn isend(&self, t: T) -> Result<(), SendError<T>> {
        match self.send(t) {
            Err(e) => Err(SendError(t)),
            Ok(_) => Ok(())
        }
    }
}
