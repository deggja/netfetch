class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters"
    homepage "https://github.com/deggja/netfetch"
    url "https://github.com/deggja/netfetch/releases/download/0.0.15/netfetch_0.0.15_darwin_amd64.tar.gz"
    sha256 "2bbf7cb3d50e477d45d6557417048a07b1d1f3884543cad3583049ef04072857"
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end