class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
    homepage "https://github.com/deggja/netfetch"
    url "https://github.com/deggja/netfetch/releases/download/0.0.62/netfetch_0.0.62_darwin_amd64.tar.gz"
    sha256 "3474896871cec3646a72d5fa9daf8b0f62656fad5100ed92406cbd186b906549"
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end