# benandjerry


To import the json into db 
cd cmd/import/
go build
./importjson -db "postgres://postgres:postgres@localhost:5433/ben-and-jerry?sslmode=disable" -file=icecream.json

RUN Application using docker-compose.sh 
./docker-compose.sh

The above command will launch the application and a webserver will be running.

To list the icecreams .

METHOD : GET 
http://localhost:8080/icecreams/list

To show a specific icecream info.

METHOD : GET 
http://localhost:8080/icecreams/list/$name

To create a new icecream info.

METHOD : POST 
http://localhost:8080/icecreams/create
BODY:  `{
	"name": "Caramel Chocolate Cheesecake1",
	"image_closed": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/caramel-chocolate-cheesecake-truffles-landing.png",
	"image_open": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/caramel-chocolate-cheesecake-truffles-landing-open.png",
	"description": "Caramel Cheesecake Ice Cream with Graham Cracker-Covered Cheesecake Truffles & Chocolate Cookie Swirls",
	"story": "In your cheesecake dreams, is it like you\u2019re spooning through a world of caramel cheesecake ice cream swirled with chocolate cookies in a wonderland of truffles filled with cheesecake? Hello? You can wake up now\u2026",
	"sourcingValues": ["Non-GMO", "Cage-Free Eggs", "Fairtrade", "Responsibly Sourced Packaging", "Caring Dairy"],
	"ingredients": ["cream", "skim milk", "water", "liquid sugar sugar", "water", "sugar", "corn syrup", "canola oil", "cream cheese pasteurized milk", "cream", "cheese cultures", "salt", "carob bean gum", "coconut oil", "egg yolks", "wheat flour", "dried cane syrup", "soybean oil", "graham flour", "eggs", "cocoa (processed with alkali)", "natural flavors", "cocoa", "guar gum", "butteroil", "milk protein concentrate", "corn starch", "salt", "soy lecithin", "tapioca starch", "pectin", "caramelized sugar syrup", "baking soda", "molasses", "honey", "carrageenan", "vanilla extract"],
	"allergy_info": "contains milk, eggs, wheat and soy",
	"dietary_certifications": "Kosher-i",
	"productId": "2191"
}`


To update an icecream info.

METHOD : POST 
http://localhost:8080/icecreams/update
BODY : same as CREATE


To delete an icecream info.

METHOD : DELETE 
http://localhost:8080/icecreams/delete/$name


