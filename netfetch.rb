class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
    homepage "https://github.com/deggja/netfetch"
    url "https://github.com/deggja/netfetch/releases/download/0.0.58/netfetch_0.0.58_darwin_amd64.tar.gz"
    sha256 "9f91c073596dd0ec511b28dcf8fc6ebe1b65e8da9310cb2d26a308befe849e65"
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end