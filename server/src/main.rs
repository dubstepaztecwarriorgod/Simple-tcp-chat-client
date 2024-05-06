use std::net::{TcpListener, TcpStream};
use std::io::{self, Read, Write};
use std::sync::{Arc, Mutex};
use std::env; 
use std::thread;
// Instead of spawining lots of threads a que could work too
// Not to sure how good of an idea this is though when many clients connect
// use std::collections::BinaryHeap;

const MAX_MESSAGE_SIZE: usize = 2024;
const DEFAULT_PORT: &str = "127.0.0.1:8080";

fn main() -> io::Result<()> {
    let address = env::args()
        .nth(1)
        .unwrap_or(DEFAULT_PORT.into());

    let listner = TcpListener::bind(&address)
        .expect(format!("Couldn't bind to the address: {}", address).as_str());

    println!("Server has been bound to address: {}", address);
 
    let clients: Arc<Mutex<Vec<TcpStream>>> = Arc::new(Mutex::new(Vec::new()));

    for connection in listner.incoming() {
        match connection {
            Ok(stream) => {
                let clients = Arc::clone(&clients);
                clients.lock().unwrap().push(stream.try_clone().unwrap());
                thread::spawn(move || {
                    if let Err(error) = handle_connection(stream, clients) {
                        eprintln!("Error: {}", error)
                    }
                });
            },
            Err(error) => eprintln!("Couldn't accept connection. Error: {}", error)
        }
    }

    Ok(())
}

fn handle_connection(mut stream: TcpStream, clients: Arc<Mutex<Vec<TcpStream>>>) -> io::Result<()> {
    let mut buffer = [0u8; MAX_MESSAGE_SIZE];
    let sender_address = stream.peer_addr().unwrap();
    println!("Connection from {} has been accepted", sender_address);
    
    loop {
        let bytes_read = stream.read(&mut buffer)?;

        for client in &mut *clients.lock().unwrap() {
            let current_client_address = client.peer_addr().unwrap();
            if current_client_address != sender_address {
                client.write(&buffer[..bytes_read])?;
            }
        }

        buffer = [0u8; MAX_MESSAGE_SIZE];
    }
}
