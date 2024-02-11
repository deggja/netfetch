class Netfetch < Formula
  desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
  homepage "https://github.com/deggja/netfetch"

  if OS.mac?
    url "https://github.com/deggja/netfetch/releases/download/3.2.12/netfetch_3.2.12_darwin_amd64.tar.gz"
    sha256 "d472c7aad2197f83edd84b601701a29097aab4d67a57f320a356f6f7b8ab2911"
  elsif OS.linux?
    url "https://github.com/deggja/netfetch/releases/download/3.2.12/netfetch_3.2.12_linux_amd64.tar.gz"
    sha256 "442f571b92376429b56de827a6e00628afbd59a20f932fc6fe566a2440ed2877"
  end

  def install
    bin.install "netfetch"
  end

  test do
    system "#{bin}/netfetch", "version"
  end
end
