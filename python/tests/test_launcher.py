import unittest

import yaju_bot


class LauncherTests(unittest.TestCase):
    def test_version(self):
        self.assertRegex(yaju_bot.__version__, r"^\d+\.\d+\.\d+")

    def test_linux_amd64_asset(self):
        self.assertEqual(
            yaju_bot._platform_asset_for("Linux", "x86_64"),
            ("yaju-bot-linux-amd64.tar.gz", "yaju-bot"),
        )

    def test_linux_arm64_asset(self):
        self.assertEqual(
            yaju_bot._platform_asset_for("Linux", "aarch64"),
            ("yaju-bot-linux-arm64.tar.gz", "yaju-bot"),
        )

    def test_windows_asset(self):
        self.assertEqual(
            yaju_bot._platform_asset_for("Windows", "AMD64"),
            ("yaju-bot-windows-amd64.zip", "yaju-bot.exe"),
        )

    def test_invalid_token(self):
        with self.assertRaises(ValueError):
            yaju_bot.run("")


if __name__ == "__main__":
    unittest.main()
