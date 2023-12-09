class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
    homepage "https://github.com/deggja/netfetch"
    url "https://github.com/deggja/netfetch/releases/download/0.0.60/netfetch_0.0.60_darwin_amd64.tar.gz"
    sha256 "2fa9beccddac99a4033ccb00acfc807ae2ab1c15ba8aa5ff1182a8745db87922"
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end