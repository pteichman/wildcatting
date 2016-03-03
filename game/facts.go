package game

var facts = []string{
	`Climate scientists are currently neutral as to whether human causes are the the major drivers of Global Warming.`,
	`Climate scientists have little faith in their models, the major tool used to "predict" climate change.`,
	`The consensus is that climate scientists' models are likely wrong, and there has been little movement from that position in the last decade.`,
	`Luis de Moscoso, a survivor of the DeSoto expedition, recorded the first sighting of oil in Texas.`,
	`The Corsicana oilfield developed gradually and peaked in 1900, when it produced more than 839,000 barrels of oil.`,
	`The first economically significant discovery of oil in Texas came in 1894 in Navarro County near Corsicana.`,
	`The legendary D. Harold (Dry Hole) Byrd was born in Detroit, Texas, on April 24, 1900, the youngest of five sons and three daughters of Edward and Mary (Easley) Byrd.`,
	`Thousands of Texans have been touched by Texas' black gold through the philanthropy of people who have made fortunes from its discovery, production and processing.`,
	`The legendary wildcatter Sid Richardson started the family oil business when his mother lent him $40 for train fare to West Texas to "put some deals together."`,
	`Perry Richardson Bass became a favorite nephew of his uncle, legendary wildcatter Sid Richardson, and decided to make his living drilling for oil.`,
	`Perry Richardson Bass, an avid fisherman, never said goodbye. His parting words were always, 'Tight lines and screaming reels.'`,
	`Perry Richardson Bass called his business a "game," adding, "We'll do anything honest to make a living."`,
	`In 1935, Perry Bass, a student at Yale University, and his uncle form Richardson & Bass, an oil venture.  After two dry holes, Richardson hits with the discovery well of the fabulously rich Keystone Field in West Texas' Winkler County.`,
	`Sid Williams Richardson was a Texas oilman, cattleman and philanthropist known for his association with the city of Fort Worth.`,
	`Sid Williams Richardson was born on 25th April, 1891, in Athens, Texas.`,
	`In 1919 Sid Richardson established his own oil company in Fort Worth.  In 1921 the oil market collapsed and he lost most of his fortune.`,
	`Sid Richardson had originally been a supporter of the Democratic Party and was associated with a group of right-wing politicians that included Dwight Eisenhower and Lyndon B. Johnson.`,
	`Sid Richardson began ranching in the 1930s and developed a love of Western art, particularly that of Frederic Remington and Charles M. Russell.`,
	`Texas includes 267,339 square miles, or 7.4% of the nation's total area.`,
	`Texas is popularly known as The Lone Star State.`,
	`Sam Houston, arguably the most famous Texan, was actually born in Virginia.`,
	`More species of bats live in Texas than in any other part of the United States.`,
	`Texas is more than an area. Texas is an idea and an experience that transcends present geographical boundaries.`,
	`Texas quickly became one of the leading oil producing states in the U.S., along with Oklahoma and California; soon the nation overtook the Russian Empire as the top producer of petroleum.`,
	`Notable wildcatters include Glenn McCarthy, Thomas Baker Slick, Sr. and Mike Benedum, Joe Trees, Clem S. Clarke, and Columbus Marion Joiner.`,
	`According to tradition, the origin of the term in the petroleum industry comes from Wildcat Hollow in Oil Creek State Park located near Titusville, Pennsylvania. Wildcat Hollow was one of the many productive fields in the early oil era.`,
	`Between 1920 and 1922, the town of Breckenridge in rural North Texas grew from about 1,500 people to nearly 30,000.`,
	`The town of Kilgore in eastern Texas grew from about 500 to 12,000 between 1930 and 1936 following the discovery of the East Texas field.`,
	`The majority of the pioneering of and searching for new oilfields in this era was done by independents, not big business interests.`,
	`An indirect effect of the boom was the growth of gambling and prostitution in many communities.`,
	`Glenn McCarthy was a modest oil worker who pioneered wells around what the Houston area. His love of bourbon led him to establish the WildCatter bourbon label.`,
	`Known as "Silver Dollar Jim", for his habit of carrying silver dollars and tossing them to doormen, the poor, and anyone that waited on him, West Jr. is regarded by many as the most flamboyant of Houston oilmen. His lavish spending habits and his proclivity for amateur law enforcement were well known.`,
	`For Texans, the 20th century did not begin on January 1, 1901, as it did for everyone else. It began nine days later, on Jan. 10, when, spurting drilling pipe, mud, gas and oil, the Lucas No. 1 well blew in at Spindletop near Beaumont.`,
	`When oil came gushing into Texas early in the 20th century, petroleum began to displace agriculture as the principal engine driving the economy of the state, and Texans' lives were even more drastically affected than they had been by railroads.`,
	`The playing out of pumped-out oil fields led to the death of any number of once-flourishing towns. Betting fortunes on what turned out to be dusters resulted in the bankruptcies of companies and individuals.`,
	`Texas oil has affected the lives of millions of Texans not directly involved in the oil business – Texans who receive neither a paycheck nor a royalty check based on petroleum.`,
	`Oil has profoundly changed the culture of the state, and it continues to affect most Texans' lives in ways that may not be obvious to the casual observer.`,
	`The presence of natural oil seeps in Texas had been known for hundreds of years before Europeans arrived in the area. Indians in Texas are said to have told European explorers that the substance had medicinal uses.`,
	`In July 1543, the remnants of Spanish explorer Hernando de Soto's expedition, led by Luis de Moscoso Alvarado, were forced ashore along the Texas coast between Sabine Pass and High Island. Moscoso reported that the group found oil floating on the surface of the water and used it to caulk their boats.`,
	`While drilling for water in 1886, Bexar County rancher George Dullnig found a small quantity of oil, but he did not attempt commercial production.`,
	`City crews in Corsicana were drilling for water in 1894, when they made the first economically significant oil discovery in Texas. That well was abandoned because the drillers needed to find water, not oil.`,
	`The oil discovery that jump-started Texas' transformation into a major petroleum producer and industrial power was Spindletop.`,
	`Spindletop triggered an influx of hundreds of eager wildcatters – including former Governor James Stephen Hogg – lusting after a piece of the action, as well as thousands of workers looking for jobs.`,
	`The Texas Oil Boom was California's fabled Gold Rush of 50 years earlier repeated on the Texas Gulf Coast with rotary drill bits and derricks instead of pick axes and gold pans.`,
	`The Texas Oil Boom turned into a feeding frenzy of human sharks: scores of speculators sniffing out a quick buck; scam artists peddling worthless leases; and prostitutes, gamblers and liquor dealers, all looking for a chunk of the workers' paychecks.`,
	`Within three years, several additional major fields were developed within a 150-mile radius of Spindletop; Sour Lake, Batson and Humble were among them.`,
	`Companies were soon established to develop the Gulf Coast oil fields. Many of them became the industry giants of today: Gulf Oil; Sun Oil Company; Magnolia Petroleum Company; the Texas Company; and Humble Oil, which later affiliated with Standard Oil of New Jersey and became Esso, then today's Exxon.`,
	`The discovery of the Spindletop oil field had an almost incalculable effect on world history, as well as Texas history.`,
	`Texas oil production was 836,039 barrels in 1900. In 1902, Spindletop alone produced more than 17 million barrels, or 94 percent of the state's production. As a result of the glut, oil prices dropped to an all-time low of 3 cents a barrel, while water in some boom towns sold for 5 cents a cup.`,
	`Between 1902 and 1910, oil fever spread through North Central Texas, with finds at Brownwood, Petrolia and Wichita Falls.`,
	`The wealth of oil at Ranger and elsewhere in the state encouraged railroads to switch their locomotives from coal to oil and helped kill the coal-mining town of Thurber.`,
	`Oil was found west of Burkburnett in Wichita County in 1912, followed by another oil field in the town itself in 1918. The feverish activity that followed inspired the 1940 movie Boom Town, starring Clark Gable, Spencer Tracy, Claudette Colbert and Hedy Lamarr.`,
	`Unexpectedly heavy traffic on the often-unpaved streets created massive clouds of dust during dry weather – dust that invaded every corner and settled on every surface. In wet weather, the streets became vehicle-swallowing mudholes.`,
	`In October 1930, the Daisy Bradford No. 3 well blew in near Turnertown and Joinerville in Rusk County, opening the East Texas field, the biggest field of all.`,
	`The biggest leasing campaign in history ensued, and the activity spread to include Kilgore, Longview and many points north. Overproduction soon followed, as oil derricks sprouted thick as bamboo all over the field. With no well-spacing regulations and no limits on production, the price of oil nosedived again.`,
	`On Aug. 17, 1931, Gov. Ross S. Sterling ordered the National Guard into the East Texas field, which he placed under martial law.`,
	`By the time the East Texas field was developed, Texas' economy was powered not by agriculture, but by petroleum.`,
	`Another change brought about by the discovery of oil was the enrichment of the state treasury after the legislature authorized an oil-production tax in 1905. The first full year the tax was collected, the public coffers swelled by $101,403.`,
	`By 1919, the revenue from the oil-production tax was more than $1 million; by 1929, it was almost $6 million.`,
	`Many thousands of students attending Texas universities have benefited from oil. The boon that they have enjoyed began with Mirabeau B. Lamar, known as the "Father of Texas Education."`,
	`Texas public schools have benefited from oil, as well. In 1839, the Congress of the Republic appropriated from the public domain three leagues of land (one league is about 4,400 acres) to each county for public schools.`,
	`In the century since Spindletop roared to life on the Texas Gulf Coast, oil has touched the lives of many Texans, and it continues to provide benefits to residents of the Lone Star State, as well as to people throughout the country.`,
	`Production from rocks of Paleozoic age occurs primarily from North Central Texas westward to New Mexico and southwestward to the Rio Grande, but there is also significant Paleozoic production in North Texas.`,
	`Mesozoic rocks are the primary hydrocarbon reservoirs of the East Texas Basin and the area south and east of the Balcones Fault Zone. Cenozoic sandstones are the main reservoirs along the Gulf Coast and offshore state waters.`,
	`Indians found oil seeping from the soils of Texas long before the first Europeans arrived. They told explorers that the fluid had medicinal values.`,
	`Melrose, in Nacogdoches County, was the site in 1866 of the first drilled well to produce oil in Texas. The driller was Lyne T. Barret. Barret used an auger, fastened to a pipe, and rotated by a cogwheel driven by a steam engine — a basic principle of rotary drilling that has been used since, although with much improvement.`,
	`Jan. 10, 1901, is the most famous date in Texas petroleum history.`,
	`The East Texas field, biggest of them all, was discovered near Turnertown and Joinerville, Rusk County, by veteran wildcatter C. M. (Dad) Joiner in October 1930.`,
	`The easy-going rural life of East Texas changed drastically with the discovery of oil in 1930 and 1931 – years of hardship, scorn, luck and wealth which brought people, ideas, institutions and national attention to East Texas.`,
	`In 1929, a 70-year-old wildcatter, Columbus Marion “Dad” Joiner, unsuccessfully drilled two dry holes south of Kilgore. Then in May, Joiner spudded a third hole on the Daisy Bradford farm in Rusk County. It was not until Oct. 3, 1930 that a production test was done, resulting in a gusher – the discovery well, Daisy Bradford No. 3.`,
}
