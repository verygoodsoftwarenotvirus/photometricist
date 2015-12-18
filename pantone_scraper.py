from selenium import webdriver

PATH_TO_WEBDRIVER = "/webdriver/chromedriver_215"
driver = webdriver.Chrome(os.getcwd() + PATH_TO_WEBDRIVER)

for x in range(10000):
    url = "http://www.pantone.com/pages/pantone/colorfinder.aspx?c_id={0}".format(x)