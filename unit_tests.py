import unittest
import color_analysis


class ColorAnalysisTest(unittest.TestCase):
    def test_color_range(self):
        red_color_range = ("#FF0000", "#C83232")
        compliant_red = "#E41919"
        non_compliant_red = "#B53232"

        self.assertTrue(color_analysis.color_is_in_range(compliant_red,
                                                         red_color_range))
        self.assertFalse(color_analysis.color_is_in_range(non_compliant_red,
                                                          red_color_range,
                                                          margin_of_error=None))