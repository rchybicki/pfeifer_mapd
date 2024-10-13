package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"

	"capnproto.org/go/capnp/v3"
	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type TmpNode struct {
	Latitude  float64
	Longitude float64
}
type TmpWay struct {
	Id                       int64  
	Name                     string
	Ref                      string
	Hazard                   string
	MaxSpeed                 float64
	MaxSpeedAdvisory float64
	MaxSpeedPractical        float64
	MaxSpeedPracticalForward float64
	MaxSpeedPracticalBackward float64
	MaxSpeedForward          float64
	MaxSpeedBackward         float64
	Lanes                    uint8
	MinLat                   float64
	MinLon                   float64
	MaxLat                   float64
	MaxLon                   float64
	OneWay                   bool
	Nodes                    []TmpNode
}

type Area struct {
	MinLat float64
	MinLon float64
	MaxLat float64
	MaxLon float64
	Ways   []TmpWay
}

var (
	GROUP_AREA_BOX_DEGREES = 2
	AREA_BOX_DEGREES       = float64(1.0 / 4) // Must be 1.0 divided by an integer number
	OVERLAP_BOX_DEGREES    = float64(0.01)
	WAYS_PER_FILE          = 2000
)

var maxSpeedOverrides = map[int64]float64{
	
  1167942324: 45, //30 Przejazd kolejowy pomiędzy Domasławiem a Bielanami

  //Wjazd do Bielan od strony Domasławia
  28345080:   65, //90 70 practical
  1172030126: 50, // 50 20 Practical pomiędzy łezkami

  //Tyniec Domasławska
  1169316187: 13, // 40 przed domem  
  1167942315: 22, // 30 hopka
  1167942314: 42, // 40
  35551085:   45, // 40
  1169316186: 40, // 40
  1167942316: 45, // 40
  193054194:  45, // 40
  1167942318: 45, // 40
  1167942322: 45, // 40
  1167942317: 30, // 30 hopka 20 practical
  1169294668: 35, // 40 practical 20 - łezka z hopką
  1169295507: 60, // 50
  1168346113: 25, // 30 //ostry zakret na koncu

  //Tyniec Szczęśliwa
  1171211140: 5, //wjazd
  1171211141: 5, //wjazd
  133979428:  7,

  //Tyniec Świdnicka
  1167942313: 35, // 30 hopka 20 practical
  913794171:  45, // 20 przy hopce rownolegla
  141340333:  45, // 20 przed hopką prostopadła
  38376718:   45, // 20 przed hopką prostopadła
  1167942312: 45, // 40
  193054196:  45, // 40
  43115448:   45, // 40
  22926964:   45, // 40
  193054197:  45, // 40
  1169807993: 42, // 40
  1169807992: 45, // 40
  185542419:  50, // 40
  941773896:  50, // 90
  941773892:  50, // 90
  941773891:  50, // 70
  941773895:  50, // 50
  941773893:  50, // 90
  941773894:  50, // 90
  186331570:  50, // 20 //Szkolna, boczna Świdnickiej
  236186142:  60, // 90

  //Domasław Tyniecka
  1167942321: 25, // 30 hopka
  1169295578: 30, // 50 przejazd kolejowy
  511534356:  40, // 50
  941773890:  55, // 70
  1169297710: 55, // 70
  253022529:  55, // 70

  //Mokronos Drogowców
  185988024: 70, // 90
  448924246: 70, // 90
  977768095: 70, // 90
  977768093: 70, // 90
  977768094: 70, // 90

  //Ślęża Wysoka
  941777387: 60, // 90

  //Ślęża przystankowa
  32723286: 40, // 50
  941777389: 40, // 50
  32723263: 55, // 50


  //Droga serwisowa przy przedszkolu
  174140502: 50, // 90
  168058572: 50, // 90

  //Ślęża nad obwodnicą
  701191713: 50, // 30
  26840348:  45, // 30
  701191714: 50, // 30

  //Wjazd na obwodnice z ronda tyniec w stronę miasta
  134429085:  65,  // 50
  360842016:  60,  // 40
  1167942327: 120, // 40

  //Zjazd z obwodnicy od strony miasta na rondo tyniec
  134429084: 60, // 70
  360842012: 60, // 70
  520977980: 50, // 50 40 practical

  //Wjazd na obwodnice z Mokronosu w stronę Tyńca
  223324845:  50,  // 50
  1169806867: 60,  // 50
  807593955:  65,  // 60
  111853026:  120, // 60
  1176894307: 120, // 60

  //Wjazd na obwodnicę z Mokronosu w stronę Miasta
  111814392: 100, // 50

  //Zjazd z obwodnicy od Tyńca w stronę Mokronosu
  111814380: 65, // 50
  223324844: 65, // 70 60 practical
  223324842: 45, // 40

  //Zjazd z obwodnicy od miasta w stronę Mokronosu
  111814378: 65, // 50

  //Wjazd na rondo w stronę Tyńca od strony obwodnicy
  112228342: 40, // 50 20 Practical

  //Wjazd na rondo w stronę Wrocławia od strony obwodnicy Mokronos
  111814391: 35, // 40 20 Practical

  //Wjazd na rondo w stronę Obwodnicy Mokronos od strony Wrocławia
  176686896:  40, // 40 30 practical
  1169300028: 35, // 40 20 Practical

  //Zjazd z obwodnicy od Kobierzyc w stronę Tyńca
  112228338: 70, // 40
  112260512: 50, // 50
  
  //A4 wjazd na obwodnicę w stronę Tyńca od strony Bielan
  272595878: 55, // 40
  272595879: 80, // 40

  //Zjazd z A4 od strony Tyńca na Bielany
  322214383: 55, // 40
  330027681: 55, // 40
  272688751: 45, // 40
  272688747: 45, // 40
  249134773: 45, // 40
  15800485:  55, // 40

  //Łacznik z A4 na obwodnicę w stronę miasta:
  39192514:  90, // 60
  286767071: 90, // 60
  112317785: 80, // 60
  286767067: 90, // 60

  //Łacznik z A4 na obwodnicę w stronę Kobierzyc:
  39192520:  90, // 60
  112228317: 80, // 60

  //Łacznik z obwodnicy na A4 w stronę Katowic
  111977997: 60, // 50
  111977827: 60, // 50
  248064919: 60, // 50
  316438185: 65, // 50
  316438184: 70, // 50​

  //Łacznik z obwodnicy na S5 w stronę Rawicza:
  121815495: 65,  // 90
  247934290: 65,  // 60
  122169242: 70,  // 60
  122169237: 70,  // 60
  122169243: 100, // 60

  //Łącznik z S5 od rawicza na obwodnicę w stronę Bielan
  388700660: 80,  // 60
  545467576: 70,  // 60
  122169239: 70,  // 60
  122169245: 65,  // 60​
  272255631: 65,  // 90
  564014164: 65,  // 70
  272256020: 65,  // 70
  272258414: 60,  // 80
  309923786: 45,  // 50
  388697213: 100, // 50

  //Za zjazdem z S5 na Rawicz w stronę Bełcza
  208286685: 60, // 70
  951950365: 60, // 70
  186474835: 60, // 70
  186474837: 45, // 50
  186474838: 60, // 70
  186474836: 60, // 70

  //Drogi pomiędzy s5 a Wąsoszem
  208287752: 80, // 90
  613473800: 80, // 90

  //Droga pomiędzy domasławiem a rondem w stronę Kobierzyc
  549996850: 60, // 90
  446503133: 60, // 90
  420299397: 60, // 90
  738968783: 50, // 90
  738968784: 50, // 90
  549996852: 50, // forward 50 backward 90
  549996853: 50, // forward 50 backward 90
  822118623: 50, // 90
  305902785: 50, // 90

  //Wąsosz Bełcz
  834278953: 70, // 90
  834278952: 60, // 90
  834278950: 70, // 90
  834278949: 70, // 90

  //pomiędzy rondami w strone Bielan od Tyńca
  520977978: 65, // 90
  185542421: 65, // 70
  119181422: 65, // 60
  118272068: 65, // 60
  119181423: 65, // 60
  904992498: 65, // 90
  904992499: 65, // 90
  206282955: 60, // 90
  546342020: 60, // 60
  546347506: 50, // 60
  20555728:  50, // 60

  //Smolec Chłopska
  133396269:  45, // 40
  1123182933: 45, // 40
  1181192632: 45, // 40
  1082811837: 45, // 40
  1181192630: 45, // 40
  1082811849: 70, // 90
  25118913:   70, // 90
  308533298:  70, // 90
  1063349750: 50, // 90

  //Mokronos Stawowa
  1171679035: 42, // 40
  25118902:   42, // 40
  25118903:   55, // 90
  174140507:  55, // 90
  448924251:  55, // 90
  
  //Mokronos Wrocławska
  1167946125: 45, // 40
  1153652251: 45, // 90 (backward), 30 (forward)
  1172221202: 45, // 40
  782197059:  55, // 60 (backward), 90 (forward)
  1172221203: 55, // 60 (backward), 90 (forward)
  1153652252: 55, // 40 (backward), 90 (forward)

  //Wrocław Zabrodzka
  25010733: 50, // 90

  //Wrocław Peronowa
  27037906: 25, // 30
  186301632: 25, // 30​

  //Wrocław Rakietowa
  549197256: 45, // 30
  306593350: 45, // 30
  53207931:  35, // 30

  //Wrocław Wyścigowa,
  847529541: 50, // 40
  164756672: 50, // 60
  546364343: 50, // 40
  307506268: 50, // 40
  492667764: 50, // 40
  513617269: 50, // 50
  492667768: 50, // 50
  492667772: 50, // 50
  492667794: 50, // 50
  546362994: 50, // 50

  //Wrocław Aleja Karkonoska obie strony przy bielanach
  194192365: 65, // 60
  194187872: 65, // 60
  124879130: 65, // 60
  330027680: 65, // 60
  330027165: 65, // 60

  //Wrocław Wiejska
  550476301: 45, // 50
  185988018: 45, // 50

  //Wrocław Aleja Karkonoska w stronę miasta
  897762977: 65, // 60
  897762976: 65, // 60
  897762975: 65, // 60
  897762974: 65, // 60
  18933198:  65, // 60
  291794288: 65, // 60
  18933205:  65, // 60
  18933202:  65, // 60
  60102381:  65, // 60
  18930495:  65, // 60
  232096918: 65, // 60
  18930510:  65, // 60
  18930504:  65, // 60
  481290473: 60, // 50
  19046871:  60, // 50
  19046874:  50, // 50
  15779094:  50, // 50
  28458096:  65, // 60
  354270946: 65, // 60

  //Wrocław Świeradowska
  313012037: 45, // 50
  186973279: 45, // 50
  
  //Wrocław Aleja Karkonoska w stronę Bielan
  28458105:  65, // 60
  16140514:  65, // 60
  353541338: 65, // 60
  307506251: 65, // 60
  331977680: 65, // 60
  331977689: 65, // 60
  482103448: 65, // 60
  307506262: 65, // 60
  31351753:  65, // 60
  307506257: 65, // 60
  307506260: 65, // 60
  307506253: 65, // 60
  186226205: 65, // 60
  28458097:  65, // 50
  307287570: 65, // 60
  307287566: 65, // 60
  307287569: 65, // 60
  307287567: 65, // 60
  92386461:  65, // 60
  92386463:  65, // 60
  492370128: 65, // 60
  492370126: 65, // 60
  492370125: 65, // 60
  92386462:  65, // 60
  307287568: 65, // 60
  18927467:  65, // 60
  18669869:  65, // 60
  18669883:  65, // 60
  161107515: 65, // 60
  18933194:  65, // 60
  794400374: 65, // 60
  291795821: 65, // 60
  897762973: 65, // 60
  18669288:  65, // 60
  60102383:  65, // 60
  194187874: 65, // 60
  897762972: 65, // 60

  //Wrocław Grabiszyńska
  322072125:  42, // 50
  15221185:   55, // 50
  1056683531: 55, // 50
  322072115:  55, // 50
  1056664607: 55, // 50
  235394664:  55, // 50

  //Wrocław Zwycięska
  258561701: 45, // 50
  967227215: 45, // 50
  967227218: 45, // 50
  679517190: 45, // 50
  679533594: 45, // 50
  679517191: 45, // 50
  679517192: 45, // 50
  679517193: 45, // 50
  679517194: 45, // 50
  679517195: 45, // 50
  967227219: 45, // 50
  679517196: 45, // 50
  823756308: 45, // 50
  333604784: 45, // 50
  679517199: 45, // 50
  679517198: 45, // 50
  679517197: 45, // 50
  508766599: 45, // 50
  508766598: 45, // 50
  968766974: 45, // 50
  968766971: 45, // 50
  968766972: 45, // 50
  679517201: 45, // 50 
  968766970: 45, // 50
  968766973: 45, // 50

  //Wrocław Radosna ołtaszyn
  1000512183: 50, // 50
  1000512182: 50, // 50
  26840374: 50, // 50
  197006758: 50, // 50
  26840373: 50, // 50
  941777386: 60, // 50
  514419194: 60, // 90
  26840350: 60, // 50


  //Wrocław Wisniowa 
  28458756:   55, // 50
  28458757:   55, // 50
  304347103:  55, // 50
  33782145:   55, // 50
  22673803:   55, // 50
  355999428:  55, // 50
  189500455:  55, // 50
  304356445:  55, // 50
  22673800:   55, // 50
  1154182991: 55, // 50
  830093165:  55, // 50
  304356131:  55, // 50
  321659947:  55, // 50
  134516705:  55, // 50
  32679156:   55, // 50
  15804616:   55, // 50

  //Wrocław Powstańców Śląskich
  28459713:  55, // 50
  353541346: 55, // 50
  353541355: 55, // 50
  353541345: 55, // 50
  222625434: 65, // 60
  353541343: 65, // 60
  353541354: 65, // 60
  353541339: 65, // 60
  353541341: 65, // 60
  353541342: 65, // 60
  353541348: 65, // 60
  28458100:  65, // 60
  353541350: 65, // 60
  685384548: 65, // 60


  //Wrocław Hallera
  22673799:   55, // 50
  321659949:  55, // 50
  321659948:  55, // 50
  28458758:   55, // 50
  304185996:  55, // 50
  353537455:  55, // 50
  28458759:   55, // 50
  366609234:  55, // 50
  366609235:  55, // 50
  304186366:  55, // 50
  304186368:  55, // 50
  353537456:  55, // 50
  353537454:  55, // 50
  353537452:  55, // 50
  353537453:  55, // 50
  361862364:  55, // 50
  28460407:   55, // 50
  28460404:   55, // 50
  361862361:  55, // 50
  304186826:  55, // 50
  28460406:   55, // 50
  353537457:  55, // 50
  28460403:   55, // 50
  1159841445: 55, // 50
  812988823:  55, // 50
  1159841446: 55, // 50
  32798115:   55, // 50
  1046584787: 55, // 50
  32798113:   55, // 50
  15270680:   55, // 50
  185948378:  55, // 50
  
  //Wrocław Armii Krajowej
  355999427: 55, // 50
  16228087:  55, // 50
  353537451: 55, // 50
  186505212: 55, // 50
  304348168: 55, // 50
  28458082:  55, // 50
  298102294: 55, // 50
  186505208: 55, // 50
  298102297: 55, // 50
  298102300: 55, // 50


  //Wrocław Mokronoska 
  650081973:  45, // 50
  1154182997: 60, // 50
  18795673:   60, // 90
  1154182996: 60, // 50
  309925361:  60, // 50
  1017147198: 60, // 50
  695817935:  60, // 50

  //Wrocław Zabrodzka
  25121542:  50, //Brak
  794401776: 50, //Brak

  //Wrocław Parkowa
  204008101:  45, // 40
  309925360:  45, // 40
  309925359:  45, // 40
  223324841:  45, // 40
  114136566:  45, // 40
  448924250:  55, // 90

  //Wrocław Ślężna w stronę miasta
  1171175010: 25, //Hopka
  800318777:  45, // 40
  22673801:   55, // 50
  800318778:  55, // 50
  304182782:  55, // 50
  304182783:  55, // 50
  322266056:  55, // 50
  304182929:  55, // 50
  190721429:  55, // 50
  546375967:  55, // 50
  546375965:  55, // 50
  546375966:  55, // 50
  546375964:  55, // 50
  190721424:  55, // 50

  //Wrocław Ślężna w stronę Bielan
  618352949:  45, // 40
  308361432:  55, // 50

  //Wrocław Kwiatkowskiego
  15921137:  60, // 50
  15921138:  60, // 50
  189327293: 55, // 50
  695817934: 55, // 50

  //Wrocław Bałtycka
  27689501:  45, // 50
  187405564: 45, // 50
  292870205: 45, // 50
  460537159: 45, // 50
  460537157: 45, // 50
  178608660: 45, // 50
  790778928: 45, // 50

  //Wrocław Reymonta
  236322555: 45, // 50
  236322552: 45, // 50
  236322553: 40, // 50 //remont
  24320103:  40, // 50 //remont
  224205260: 40, // 50 //remont

  //Wrocław Gajowicka
  504541555: 40, // 50
  173689459: 40, // 50
  173689455: 40, // 50
  504403344: 40, // 50
  504403343: 40, // 50
  191144133: 40, // 50

  //Wrocław Tyniecka
  24983229:  55,  // 50

  //Wrocław Jeziorańskiego
  231316713: 65,  // 50
  307080161: 65,  // 50
  83410375:  65,  // 50
  370404883: 65,  // 50
  307080157: 60,  // 50

  //Wrocław Aleja Pracy
  16768472: 40,  // 50


  //Wrocław Aleja Piastów
  //Hopki
  1172811812: 25, // 30
  1177217053: 25, // 30
  1172811814: 25, // 30
  1172811854: 25, // 30
  1172811816: 25, // 30
  1172811852: 25, // 30
  1172811818: 25, // 30
  1172811850: 25, // 30
  1172811820: 25, // 30
  1172811848: 25, // 30
  1172811822: 25, // 30
  1172811824: 25, // 30
  1172811846: 25, // 30
  1172811826: 25, // 30
  1172811844: 25, // 30
  1172811828: 25, // 30
  1172811842: 25, // 30
  1172811840: 25, // 30
  1172811830: 25, // 30
  1172811838: 25, // 30
  1172811832: 25, // 30
  1172811836: 25, // 30
  1172811834: 25, // 30

  15118987:   35, // 30
  1153652257: 35, // 30
  186153784:  35, // 30
  443761216:  35, // 30
  1154182994: 35, // 30
  443761965:  35, // 30
  1154182995: 35, // 30
  437164390:  35, // 30
  1172811811: 35, // 30
  437164391:  35, // 30
  423682103:  35, // 30
  401158323:  35, // 30
  185991363:  35, // 30
  32798101:   35, // 30
  443826986:  35, // 30
  443826987:  35, // 30
  1172811833: 35, // 30
  1172811835: 35, // 30
  1177217052: 35, // 30
  1172811837: 35, // 30
  1172811839: 35, // 30
  1172811841: 35, // 30
  1172811843: 35, // 30
  1172811845: 35, // 30
  1172811847: 35, // 30
  1172811849: 35, // 30
  1172811851: 35, // 30
  1172811853: 35, // 30
  450113197:  50, // 30
  423682102:  50, // 50
  640962322:  50, // 50
  313926987:  50, // 50
  783156556:  50, // 50
  640961668:  50, // 50
  373946571:  50, // 50
  854827764:  50, // 50 backward, 40 forward
  854827765:  50, // 40 backward, 50 forward
  854827766:  50, // 50
  177522534:  55, // 50

  //Pomiędzy Górą a rondem
  151875999: 80, // 90
  94941811:  80, // 90
  438193041: 80, // 90
}

func GetBaseOpPath() string {
	exists, err := Exists("/data/media/0")
	logde(err)
	if exists {
		return "/data/media/0/osm"
	} else {
		return "."
	}
}

var BOUNDS_DIR = fmt.Sprintf("%s/offline", GetBaseOpPath())

func EnsureOfflineMapsDirectories() {
	err := os.MkdirAll(BOUNDS_DIR, 0o775)
	logwe(err)
}

// Creates a file for a specific bounding box
func GenerateBoundsFileName(minLat float64, minLon float64, maxLat float64, maxLon float64) string {
	group_lat_directory := int(math.Floor(minLat/float64(GROUP_AREA_BOX_DEGREES))) * GROUP_AREA_BOX_DEGREES
	group_lon_directory := int(math.Floor(minLon/float64(GROUP_AREA_BOX_DEGREES))) * GROUP_AREA_BOX_DEGREES
	dir := fmt.Sprintf("%s/%d/%d", BOUNDS_DIR, group_lat_directory, group_lon_directory)
	return fmt.Sprintf("%s/%f_%f_%f_%f", dir, minLat, minLon, maxLat, maxLon)
}

// Creates a file for a specific bounding box
func CreateBoundsDir(minLat float64, minLon float64, maxLat float64, maxLon float64) error {
	group_lat_directory := int(math.Floor(minLat/float64(GROUP_AREA_BOX_DEGREES))) * GROUP_AREA_BOX_DEGREES
	group_lon_directory := int(math.Floor(minLon/float64(GROUP_AREA_BOX_DEGREES))) * GROUP_AREA_BOX_DEGREES
	dir := fmt.Sprintf("%s/%d/%d", BOUNDS_DIR, group_lat_directory, group_lon_directory)
	err := os.MkdirAll(dir, 0o775)
	return errors.Wrap(err, "could not create bounds directory")
}

// Checks if two bounding boxes intersect
func Overlapping(axMin float64, ayMin float64, axMax float64, ayMax float64, bxMin float64, byMin float64, bxMax float64, byMax float64) bool {
	intersect := !(axMin > bxMax || axMax < bxMin || ayMin > byMax || ayMax < byMin)
	aMinInside := PointInBox(axMin, ayMin, bxMin, byMin, bxMax, byMax)
	bMinInside := PointInBox(bxMin, byMin, axMin, ayMin, axMax, ayMax)
	aMaxInside := PointInBox(axMax, ayMax, bxMin, byMin, bxMax, byMax)
	bMaxInside := PointInBox(bxMax, byMax, axMin, ayMin, axMax, ayMax)
	return intersect || aMinInside || bMinInside || aMaxInside || bMaxInside
}

// Generates bounding boxes for storing ways
func GenerateAreas() []Area {
	areas := make([]Area, int((361/AREA_BOX_DEGREES)*(181/AREA_BOX_DEGREES)))
	index := 0
	for i := float64(-90); i < 90; i += AREA_BOX_DEGREES {
		for j := float64(-180); j < 180; j += AREA_BOX_DEGREES {
			a := &areas[index]
			a.MinLat = i
			a.MinLon = j
			a.MaxLat = i + AREA_BOX_DEGREES
			a.MaxLon = j + AREA_BOX_DEGREES
			index += 1
		}
	}
	return areas
}

func GenerateOffline(minGenLat int, minGenLon int, maxGenLat int, maxGenLon int, generateEmptyFiles bool) {
	log.Info().Msg("Generating Offline Map")
	EnsureOfflineMapsDirectories()
	file, err := os.Open("./map.osm.pbf")
	check(errors.Wrap(err, "could not open map pbf file"))
	defer file.Close()

	// The third parameter is the number of parallel decoders to use.
	scanner := osmpbf.New(context.Background(), file, runtime.GOMAXPROCS(-1))
	scanner.SkipRelations = true
	defer scanner.Close()

	scannedWays := []TmpWay{}
	areas := GenerateAreas()
	index := 0
	allMinLat := float64(90)
	allMinLon := float64(180)
	allMaxLat := float64(-90)
	allMaxLon := float64(-180)

	log.Info().Msg("Scanning Ways")
	for scanner.Scan() {
		var way *osm.Way
		switch o := scanner.Object(); o.(type) {
		case *osm.Way:
			way = o.(*osm.Way)
		default:
			way = nil
		}
		if way != nil && len(way.Nodes) > 1 {
			tags := way.TagMap()
			lanes, _ := strconv.ParseUint(tags["lanes"], 10, 8)
			tmpWay := TmpWay{
				Id:                       int64(way.ID),
				Nodes:                    make([]TmpNode, len(way.Nodes)),
				Name:                     tags["name"],
				Ref:                      tags["ref"],
				Hazard:                   tags["hazard"],
				MaxSpeed:                 ParseMaxSpeed(tags["maxspeed"]),
				MaxSpeedAdvisory: ParseMaxSpeed(tags["maxspeed:advisory"]),
				MaxSpeedPractical:        ParseMaxSpeed(tags["maxspeed:practical"]),
				MaxSpeedPracticalForward: ParseMaxSpeed(tags["maxspeed:practical:forward"]),
				MaxSpeedPracticalBackward: ParseMaxSpeed(tags["maxspeed:practical:backward"]),
				MaxSpeedForward:          ParseMaxSpeed(tags["maxspeed:forward"]),
				MaxSpeedBackward:         ParseMaxSpeed(tags["maxspeed:backward"]),
				Lanes:                    uint8(lanes),
				OneWay:                   tags["oneway"] == "yes",
			}
			index++

			// Apply override if it exists, only for MaxSpeedPractical
			if override, exists := maxSpeedOverrides[tmpWay.Id]; exists {
				tmpWay.MaxSpeedPractical = 0.277778 * override
			}

			minLat := float64(90)
			minLon := float64(180)
			maxLat := float64(-90)
			maxLon := float64(-180)
			for i, n := range way.Nodes {
				if n.Lat < minLat {
					minLat = n.Lat
				}
				if n.Lon < minLon {
					minLon = n.Lon
				}
				if n.Lat > maxLat {
					maxLat = n.Lat
				}
				if n.Lon > maxLon {
					maxLon = n.Lon
				}
				tmpWay.Nodes[i].Latitude = n.Lat
				tmpWay.Nodes[i].Longitude = n.Lon
			}
			tmpWay.MinLat = minLat
			tmpWay.MinLon = minLon
			tmpWay.MaxLat = maxLat
			tmpWay.MaxLon = maxLon
			if minLat < allMinLat {
				allMinLat = minLat
			}
			if minLon < allMinLon {
				allMinLon = minLon
			}
			if maxLat > allMaxLat {
				allMaxLat = maxLat
			}
			if maxLon > allMaxLon {
				allMaxLon = maxLon
			}
			scannedWays = append(scannedWays, tmpWay)
		}
	}

	log.Info().Msg("Finding Bounds")
	for _, area := range areas {
		if area.MinLat < float64(minGenLat)-OVERLAP_BOX_DEGREES || area.MinLon < float64(minGenLon)-OVERLAP_BOX_DEGREES || area.MaxLat > float64(maxGenLat)+OVERLAP_BOX_DEGREES || area.MaxLon > float64(maxGenLon)+OVERLAP_BOX_DEGREES {
			continue
		}

		haveWays := Overlapping(allMinLat, allMinLon, allMaxLat, allMaxLon, area.MinLat-OVERLAP_BOX_DEGREES, area.MinLon-OVERLAP_BOX_DEGREES, area.MaxLat+OVERLAP_BOX_DEGREES, area.MaxLon+OVERLAP_BOX_DEGREES)
		if !haveWays && !generateEmptyFiles {
			continue
		}

		arena := capnp.MultiSegment([][]byte{})
		msg, seg, err := capnp.NewMessage(arena)
		check(errors.Wrap(err, "could not create capnp arena for offline data"))
		rootOffline, err := NewRootOffline(seg)
		check(errors.Wrap(err, "could not create capnp offline root"))

		for _, way := range scannedWays {
			overlaps := Overlapping(way.MinLat, way.MinLon, way.MaxLat, way.MaxLon, area.MinLat-OVERLAP_BOX_DEGREES, area.MinLon-OVERLAP_BOX_DEGREES, area.MaxLat+OVERLAP_BOX_DEGREES, area.MaxLon+OVERLAP_BOX_DEGREES)
			if overlaps {
				area.Ways = append(area.Ways, way)
			}
		}

		log.Info().Msg("Writing Area")
		ways, err := rootOffline.NewWays(int32(len(area.Ways)))
		check(errors.Wrap(err, "could not create ways in offline data"))
		rootOffline.SetMinLat(area.MinLat)
		rootOffline.SetMinLon(area.MinLon)
		rootOffline.SetMaxLat(area.MaxLat)
		rootOffline.SetMaxLon(area.MaxLon)
		rootOffline.SetOverlap(OVERLAP_BOX_DEGREES)
		for i, way := range area.Ways {
			w := ways.At(i)
			w.SetId(way.Id)
			w.SetMinLat(way.MinLat)
			w.SetMinLon(way.MinLon)
			w.SetMaxLat(way.MaxLat)
			w.SetMaxLon(way.MaxLon)
			err := w.SetName(way.Name)
			check(errors.Wrap(err, "could not set way name"))
			err = w.SetRef(way.Ref)
			check(errors.Wrap(err, "could not set way ref"))
			err = w.SetHazard(way.Hazard)
			check(errors.Wrap(err, "could not set way hazard"))
			w.SetMaxSpeed(way.MaxSpeed)
			w.SetAdvisorySpeed(way.MaxSpeedAdvisory)
			w.SetMaxSpeedPractical(way.MaxSpeedPractical)
			w.SetMaxSpeedPracticalForward(way.MaxSpeedPracticalForward)
			w.SetMaxSpeedPracticalBackward(way.MaxSpeedPracticalBackward)
			w.SetMaxSpeedForward(way.MaxSpeedForward)
			w.SetMaxSpeedBackward(way.MaxSpeedBackward)
			w.SetLanes(way.Lanes)
			w.SetOneWay(way.OneWay)
			nodes, err := w.NewNodes(int32(len(way.Nodes)))
			check(errors.Wrap(err, "could not create way nodes"))
			for j, node := range way.Nodes {
				n := nodes.At(j)
				n.SetLatitude(node.Latitude)
				n.SetLongitude(node.Longitude)
			}
		}

		data, err := msg.MarshalPacked()
		check(errors.Wrap(err, "could not marshal offline data"))
		err = CreateBoundsDir(area.MinLat, area.MinLon, area.MaxLat, area.MaxLon)
		check(errors.Wrap(err, "could not create directory for bounds file"))
		err = os.WriteFile(GenerateBoundsFileName(area.MinLat, area.MinLon, area.MaxLat, area.MaxLon), data, 0o644)
		check(errors.Wrap(err, "could not write offline data to file"))
	}
	f, err := os.Open(BOUNDS_DIR)
	check(errors.Wrap(err, "could not open bounds directory"))
	err = f.Sync()
	check(errors.Wrap(err, "could not fsync bounds directory"))
	err = f.Close()
	check(errors.Wrap(err, "could not close bounds directory"))

	log.Info().Msg("Done Generating Offline Map")
}

func PointInBox(ax float64, ay float64, bxMin float64, byMin float64, bxMax float64, byMax float64) bool {
	return ax > bxMin && ax < bxMax && ay > byMin && ay < byMax
}

var AREAS = GenerateAreas()

func FindWaysAroundLocation(lat float64, lon float64) ([]byte, error) {
	for _, area := range AREAS {
		inBox := PointInBox(lat, lon, area.MinLat, area.MinLon, area.MaxLat, area.MaxLon)
		if inBox {
			boundsName := GenerateBoundsFileName(area.MinLat, area.MinLon, area.MaxLat, area.MaxLon)
			log.Info().Str("filename", boundsName).Msg("Loading bounds file")
			data, err := os.ReadFile(boundsName)
			return data, errors.Wrap(err, "could not read current offline data file")
		}
	}
	return []uint8{}, nil
}