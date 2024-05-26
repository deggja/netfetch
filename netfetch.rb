class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"

  if OS.mac?
    url "https://github.com/deggja/netfetch/releases/download/5.2.5/netfetch_5.2.5_darwin_amd64.tar.gz"
    sha256 "056daed4bef3da2149d2bf9c1a9dc181fad2398b2c7add462d0b062090705c44"
  elsif OS.linux?
    url "https://github.com/deggja/netfetch/releases/download/5.2.5/netfetch_5.2.5_linux_amd64.tar.gz"
    sha256 "cdec364c59d5ae41a7d755d5b3c9afbccd4bd935c3b539dd699eb82b8e6bb8e0"
  end

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "version"
  end
end
