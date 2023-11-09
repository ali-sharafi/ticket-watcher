package app

import (
	"encoding/json"
	"ticket-watcher/domain"
	"ticket-watcher/pkg/utils"
)

func GetProvinces() (provinces []domain.Province) {
	jsonData := `[
    {
        "name": "تهران",
        "code": "THR"
    },
    {
        "name": "اهواز",
        "code": "AWZ"
    },
    {
        "name": "شیراز",
        "code": "SYZ"
    },
    {
        "name": "مشهد",
        "code": "MHD"
    },
    {
        "name": "بندر عباس",
        "code": "BND"
    },
    {
        "name": "اصفهان",
        "code": "IFN"
    },
    {
        "name": "تبریز",
        "code": "TBZ"
    },
    {
        "name": "کیش",
        "code": "KIH"
    },
    {
        "name": "آبادان",
        "code": "ABD"
    },
    {
        "name": "اراک",
        "code": "AJK"
    },
    {
        "name": "اردبیل",
        "code": "ADU"
    },
    {
        "name": "ارومیه",
        "code": "OMH"
    },
    {
        "name": "امیدیه",
        "code": "AKW"
    },
    {
        "name": "ایرانشهر",
        "code": "IHR"
    },
    {
        "name": "ایلام",
        "code": "IIL"
    },
    {
        "name": "بجنورد",
        "code": "BJB"
    },
    {
        "name": "بم",
        "code": "BXR"
    },
    {
        "name": "بندر لنگه",
        "code": "BDH"
    },
    {
        "name": "بوشهر",
        "code": "BUZ"
    },
    {
        "name": "بیرجند",
        "code": "XBJ"
    },
    {
        "name": "پارس آباد",
        "code": "PFQ"
    },
    {
        "name": "جاسک",
        "code": "JSK"
    },
    {
        "name": "جهرم",
        "code": "JAR"
    },
    {
        "name": "جیرفت",
        "code": "JYR"
    },
    {
        "name": "چابهار",
        "code": "ZBR"
    },
    {
        "name": "خارک",
        "code": "KHK"
    },
    {
        "name": "خرم آباد",
        "code": "KHD"
    },
    {
        "name": "خوی",
        "code": "KHY"
    },
    {
        "name": "دزفول",
        "code": "DEF"
    },
    {
        "name": "رامسر",
        "code": "RZR"
    },
    {
        "name": "رشت",
        "code": "RAS"
    },
    {
        "name": "رفسنجان",
        "code": "RJN"
    },
    {
        "name": "زابل",
        "code": "ACZ"
    },
    {
        "name": "زاهدان",
        "code": "ZAH"
    },
    {
        "name": "زنجان",
        "code": "JWN"
    },
    {
        "name": "ساری",
        "code": "SRY"
    },
    {
        "name": "سبزوار",
        "code": "AFZ"
    },
    {
        "name": "سمنان",
        "code": "SNX"
    },
    {
        "name": "سنندج",
        "code": "SDG"
    },
    {
        "name": "مراغه",
        "code": "ACP"
    },
    {
        "name": "سیرجان",
        "code": "SYJ"
    },
    {
        "name": "شاهرود",
        "code": "RUD"
    },
    {
        "name": "شهرکرد",
        "code": "CQD"
    },
    {
        "name": "طبس",
        "code": "TCX"
    },
    {
        "name": "عسلویه",
        "code": "PGU"
    },
    {
        "name": "قشم",
        "code": "GSM"
    },
    {
        "name": "کاشان",
        "code": "KKS"
    },
    {
        "name": "کرج",
        "code": "PYK"
    },
    {
        "name": "کرمان",
        "code": "KER"
    },
    {
        "name": "کرمانشاه",
        "code": "KSH"
    },
    {
        "name": "کلاله",
        "code": "KLM"
    },
    {
        "name": "گچساران",
        "code": "GCH"
    },
    {
        "name": "گرگان",
        "code": "GBT"
    },
    {
        "name": "لار",
        "code": "LRR"
    },
    {
        "name": "لامرد",
        "code": "LFM"
    },
    {
        "name": "ماکو",
        "code": "MAC"
    },
    {
        "name": "ماهشهر",
        "code": "MRX"
    },
    {
        "name": "نوشهر",
        "code": "NSH"
    },
    {
        "name": "همدان",
        "code": "HDM"
    },
    {
        "name": "یاسوج",
        "code": "YES"
    },
    {
        "name": "یزد",
        "code": "AZD"
    }
]`

	if err := json.Unmarshal([]byte(jsonData), &provinces); err != nil {
		utils.Logger.Error("Error in Unmarshal provinces:", err)
	}
	return
}
