defmodule UdpClientServer do
  @moduledoc """
  Documentation for UdpClientServer.
  """

  @default_server_port 8000
  @server_addr "2001:412:abcd:2:0013:A200:4147:8C2B"

  def launch_server do
    launch_server(@default_server_port)
  end

  def launch_server(port) do
    IO.puts "Launching server on localhost on port #{port}"
    server = Socket.UDP.open!(port, [{:version, 6}])
    serve(server)
  end

  def serve(server) do
    {data, client} = server |> Socket.Datagram.recv!
    IO.puts "Received: #{data}, from #{inspect(client)}"

    serve(server)
  end

  @doc """
  Sends `data` to the `to` value, where `to` is a tuple of 
  { host, port } like {{127, 0, 0, 1}, 1337}
  """
  def send_data(data, to) do
    sock = Socket.UDP.open!(0, [{:version, 6}])
    Socket.Datagram.send!(sock, data, to)
  end

  def send_data(i) do
    addr = Socket.Address.parse(@server_addr)
    msg = "HELLO #{i}"
    send_data(msg, {addr, @default_server_port})
    :timer.sleep(1000)
    send_data(i + 1)
  end

  def send_data do
    send_data(0)
  end
end
