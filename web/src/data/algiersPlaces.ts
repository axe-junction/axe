export interface AlgiersPlace {
  id: string;
  name: string;
  arabicName?: string;
  category: string;
  subCategory?: string;
  lat: number;
  lng: number;
  district: string;
  keywords: string[];
  popularity: number; // 1-10 scale
}

export const algiersPlaces: AlgiersPlace[] = [
  // Districts & Areas
  {
    id: "alger-centre",
    name: "Alger Centre",
    arabicName: "الجزائر الوسطى",
    category: "District",
    lat: 36.7538,
    lng: 3.0588,
    district: "Alger Centre",
    keywords: [
      "centre",
      "downtown",
      "city center",
      "وسط المدينة",
      "center",
      "alger",
      "metro hub",
      "tram central",
      "business district",
    ],
    popularity: 10,
  },
  {
    id: "bab-ezzouar",
    name: "Bab Ezzouar",
    arabicName: "باب الزوار",
    category: "District",
    lat: 36.7167,
    lng: 3.1833,
    district: "Bab Ezzouar",
    keywords: [
      "bab ezzouar",
      "باب الزوار",
      "business district",
      "bab",
      "ezzouar",
    ],
    popularity: 9,
  },
  {
    id: "hydra",
    name: "Hydra",
    arabicName: "هيدرا",
    category: "District",
    lat: 36.7669,
    lng: 3.0347,
    district: "Hydra",
    keywords: ["hydra", "هيدرا", "diplomatic", "embassy"],
    popularity: 8,
  },
  {
    id: "kouba",
    name: "Kouba",
    arabicName: "القبة",
    category: "District",
    lat: 36.7333,
    lng: 3.0833,
    district: "Kouba",
    keywords: ["kouba", "القبة", "qouba"],
    popularity: 7,
  },
  {
    id: "birtouta",
    name: "Birtouta",
    arabicName: "بئر توتة",
    category: "District",
    lat: 36.6167,
    lng: 3.0167,
    district: "Birtouta",
    keywords: ["birtouta", "بئر توتة", "bir", "touta"],
    popularity: 6,
  },
  {
    id: "cheraga",
    name: "Cheraga",
    arabicName: "شراقة",
    category: "District",
    lat: 36.6139,
    lng: 2.9444,
    district: "Cheraga",
    keywords: ["cheraga", "شراقة", "cherraga"],
    popularity: 7,
  },
  {
    id: "el-harrach",
    name: "El Harrach",
    arabicName: "الحراش",
    category: "District",
    lat: 36.7167,
    lng: 3.1333,
    district: "El Harrach",
    keywords: [
      "harrach",
      "الحراش",
      "el harrach",
      "metro station",
      "industrial",
      "metro line 1",
      "terminus",
      "bus hub",
    ],
    popularity: 8, // Increased popularity since it's a major metro stop
  },
  {
    id: "dar-el-beida",
    name: "Dar El Beida",
    arabicName: "دار البيضاء",
    category: "District",
    lat: 36.7081,
    lng: 3.2122,
    district: "Dar El Beida",
    keywords: [
      "dar el beida",
      "دار البيضاء",
      "dar",
      "beida",
      "white house",
      "الدار البيضاء",
    ],
    popularity: 7,
  },
  {
    id: "el-biar",
    name: "El Biar",
    arabicName: "البيار",
    category: "District",
    lat: 36.7653,
    lng: 3.0333,
    district: "El Biar",
    keywords: ["el biar", "البيار", "biar"],
    popularity: 6,
  },
  {
    id: "ben-aknoun",
    name: "Ben Aknoun",
    arabicName: "بن عكنون",
    category: "District",
    lat: 36.7167,
    lng: 2.9833,
    district: "Ben Aknoun",
    keywords: ["ben aknoun", "بن عكنون", "aknoun"],
    popularity: 6,
  },
  {
    id: "dely-brahim",
    name: "Dely Brahim",
    arabicName: "دالي إبراهيم",
    category: "District",
    lat: 36.7381,
    lng: 2.9978,
    district: "Dely Brahim",
    keywords: ["dely brahim", "دالي إبراهيم", "dely", "brahim"],
    popularity: 6,
  },
  {
    id: "douera",
    name: "Douera",
    arabicName: "الدويرة",
    category: "District",
    lat: 36.6767,
    lng: 2.9444,
    district: "Douera",
    keywords: ["douera", "الدويرة", "douira"],
    popularity: 5,
  },
  {
    id: "hussein-dey",
    name: "Hussein Dey",
    arabicName: "حسين داي",
    category: "District",
    lat: 36.7333,
    lng: 3.1,
    district: "Hussein Dey",
    keywords: ["hussein dey", "حسين داي", "hussain", "dey"],
    popularity: 6,
  },

  // Transportation Hubs
  {
    id: "houari-boumediene-airport",
    name: "Houari Boumediene Airport",
    arabicName: "مطار هواري بومدين",
    category: "Airport",
    subCategory: "International",
    lat: 36.691,
    lng: 3.2154,
    district: "Dar El Beida",
    keywords: [
      "airport",
      "مطار",
      "houari boumediene",
      "ALG",
      "international",
      "flights",
    ],
    popularity: 10,
  },
  {
    id: "algiers-port",
    name: "Algiers Port",
    arabicName: "ميناء الجزائر",
    category: "Port",
    lat: 36.7697,
    lng: 3.0611,
    district: "Alger Centre",
    keywords: ["port", "ميناء", "harbor", "ferry", "maritime"],
    popularity: 8,
  },
  {
    id: "tafourah-station",
    name: "Tafourah Metro Station",
    arabicName: "محطة تفورة",
    category: "Metro Station",
    lat: 36.7403,
    lng: 3.0508,
    district: "Alger Centre",
    keywords: ["tafourah", "metro", "station", "تفورة", "subway"],
    popularity: 9,
  },
  {
    id: "1er-mai-station",
    name: "1er Mai Metro Station",
    arabicName: "محطة الفاتح من مايو",
    category: "Metro Station",
    lat: 36.7542,
    lng: 3.0431,
    district: "Alger Centre",
    keywords: ["1er mai", "metro", "station", "may 1st"],
    popularity: 8,
  },

  // Universities & Education
  {
    id: "university-algiers",
    name: "University of Algiers",
    arabicName: "جامعة الجزائر",
    category: "University",
    lat: 36.7167,
    lng: 3.1833,
    district: "Bab Ezzouar",
    keywords: ["university", "جامعة", "education", "students", "univ"],
    popularity: 8,
  },
  {
    id: "usthb",
    name: "USTHB",
    arabicName: "جامعة العلوم والتكنولوجيا",
    category: "University",
    lat: 36.7081,
    lng: 3.1622,
    district: "Bab Ezzouar",
    keywords: [
      "usthb",
      "science",
      "technology",
      "university",
      "houari boumediene",
    ],
    popularity: 8,
  },
  {
    id: "enp-algiers",
    name: "École Nationale Polytechnique",
    arabicName: "المدرسة الوطنية المتعددة التقنيات",
    category: "University",
    lat: 36.7194,
    lng: 3.1756,
    district: "El Harrach",
    keywords: ["enp", "polytechnique", "engineering", "école", "national"],
    popularity: 7,
  },
  {
    id: "university-algiers-3",
    name: "University of Algiers 3",
    arabicName: "جامعة الجزائر 3",
    category: "University",
    lat: 36.7333,
    lng: 3.0833,
    district: "Dely Brahim",
    keywords: ["university", "algiers 3", "جامعة الجزائر", "ibrahim sultan"],
    popularity: 7,
  },
  {
    id: "esi-algiers",
    name: "ESI Algiers",
    arabicName: "المدرسة العليا للإعلام الآلي",
    category: "University",
    lat: 36.7167,
    lng: 3.1833,
    district: "Bab Ezzouar",
    keywords: ["esi", "computer science", "informatics", "école supérieure"],
    popularity: 7,
  },

  // Landmarks & Monuments
  {
    id: "maqam-echahid",
    name: "Maqam Echahid",
    arabicName: "مقام الشهيد",
    category: "Monument",
    lat: 36.7525,
    lng: 3.0519,
    district: "Alger Centre",
    keywords: [
      "maqam",
      "monument",
      "martyrs",
      "memorial",
      "مقام الشهيد",
      "shahid",
    ],
    popularity: 9,
  },
  {
    id: "grande-poste",
    name: "Grande Poste",
    arabicName: "البريد المركزي",
    category: "Landmark",
    lat: 36.7697,
    lng: 3.0597,
    district: "Alger Centre",
    keywords: ["grande poste", "post office", "البريد", "postal", "historic"],
    popularity: 8,
  },
  {
    id: "casbah",
    name: "Casbah of Algiers",
    arabicName: "قصبة الجزائر",
    category: "Historic Site",
    subCategory: "UNESCO World Heritage",
    lat: 36.7833,
    lng: 3.0597,
    district: "Casbah",
    keywords: ["casbah", "قصبة", "old city", "unesco", "medina", "historic"],
    popularity: 10,
  },
  {
    id: "ketchaoua-mosque",
    name: "Ketchaoua Mosque",
    arabicName: "جامع كتشاوة",
    category: "Mosque",
    lat: 36.7831,
    lng: 3.0598,
    district: "Casbah",
    keywords: ["ketchaoua", "mosque", "كتشاوة", "جامع", "islamic"],
    popularity: 7,
  },

  // Shopping & Commercial
  {
    id: "riadh-el-feth",
    name: "Riadh El Feth",
    arabicName: "رياض الفتح",
    category: "Shopping Center",
    lat: 36.7403,
    lng: 3.0508,
    district: "Alger Centre",
    keywords: ["riadh el feth", "shopping", "mall", "رياض الفتح", "commercial"],
    popularity: 8,
  },
  {
    id: "centre-commercial-bab-ezzouar",
    name: "Centre Commercial Bab Ezzouar",
    arabicName: "المركز التجاري باب الزوار",
    category: "Shopping Center",
    lat: 36.7156,
    lng: 3.1844,
    district: "Bab Ezzouar",
    keywords: ["shopping", "mall", "commercial center", "bab ezzouar"],
    popularity: 7,
  },
  {
    id: "marche-nelson-mandela",
    name: "Marché Nelson Mandela",
    arabicName: "سوق نيلسون مانديلا",
    category: "Market",
    lat: 36.7538,
    lng: 3.0588,
    district: "Alger Centre",
    keywords: ["market", "marché", "سوق", "nelson mandela", "shopping"],
    popularity: 6,
  },
  {
    id: "ardis",
    name: "Ardis Shopping Center",
    arabicName: "مركز أرديس التجاري",
    category: "Shopping Center",
    lat: 36.7167,
    lng: 2.9833,
    district: "Ben Aknoun",
    keywords: ["ardis", "shopping", "center", "mall"],
    popularity: 7,
  },

  // Hospitals & Healthcare
  {
    id: "chu-mustapha-pacha",
    name: "CHU Mustapha Pacha",
    arabicName: "مستشفى مصطفى باشا",
    category: "Hospital",
    lat: 36.7597,
    lng: 3.0453,
    district: "Alger Centre",
    keywords: ["hospital", "chu", "mustapha pacha", "مستشفى", "medical"],
    popularity: 8,
  },
  {
    id: "hopital-zmirli",
    name: "Hôpital Zmirli",
    arabicName: "مستشفى زميرلي",
    category: "Hospital",
    lat: 36.7167,
    lng: 3.1833,
    district: "El Harrach",
    keywords: ["hospital", "zmirli", "مستشفى زميرلي", "medical"],
    popularity: 7,
  },

  // Hotels & Accommodation
  {
    id: "hotel-aurassi",
    name: "Hotel Aurassi",
    arabicName: "فندق الأوراسي",
    category: "Hotel",
    subCategory: "5 Star",
    lat: 36.7472,
    lng: 3.0569,
    district: "Alger Centre",
    keywords: ["hotel", "aurassi", "فندق الأوراسي", "luxury", "5 star"],
    popularity: 8,
  },
  {
    id: "sheraton-club-des-pins",
    name: "Sheraton Club des Pins",
    arabicName: "شيراتون نادي الصنوبر",
    category: "Hotel",
    subCategory: "5 Star",
    lat: 36.7833,
    lng: 2.9167,
    district: "Sidi Fredj",
    keywords: ["sheraton", "club des pins", "resort", "luxury"],
    popularity: 7,
  },

  // Parks & Recreation
  {
    id: "jardin-dessai",
    name: "Jardin d'Essai",
    arabicName: "حديقة التجارب",
    category: "Park",
    lat: 36.7167,
    lng: 3.1167,
    district: "El Hamma",
    keywords: ["jardin", "park", "حديقة", "botanical", "garden", "essai"],
    popularity: 8,
  },
  {
    id: "parc-ben-aknoun",
    name: "Parc Ben Aknoun",
    arabicName: "حديقة بن عكنون",
    category: "Park",
    lat: 36.7167,
    lng: 2.9833,
    district: "Ben Aknoun",
    keywords: ["parc", "park", "ben aknoun", "حديقة", "recreation"],
    popularity: 7,
  },

  // Beaches & Coastal Areas
  {
    id: "sablettes-beach",
    name: "Sablettes Beach",
    arabicName: "شاطئ السابليت",
    category: "Beach",
    lat: 36.7667,
    lng: 3.0333,
    district: "Alger Centre",
    keywords: ["beach", "sablettes", "شاطئ", "swimming", "coast"],
    popularity: 8,
  },
  {
    id: "club-des-pins-beach",
    name: "Club des Pins Beach",
    arabicName: "شاطئ نادي الصنوبر",
    category: "Beach",
    lat: 36.7833,
    lng: 2.9167,
    district: "Sidi Fredj",
    keywords: ["beach", "club des pins", "resort", "شاطئ", "luxury"],
    popularity: 9,
  },

  // Museums & Culture
  {
    id: "bardo-museum",
    name: "Bardo National Museum",
    arabicName: "المتحف الوطني باردو",
    category: "Museum",
    lat: 36.7528,
    lng: 3.0494,
    district: "Alger Centre",
    keywords: ["museum", "bardo", "متحف", "culture", "national", "art"],
    popularity: 7,
  },
  {
    id: "mujahideen-museum",
    name: "Mujahideen Museum",
    arabicName: "متحف المجاهد",
    category: "Museum",
    lat: 36.7528,
    lng: 3.0494,
    district: "Alger Centre",
    keywords: [
      "museum",
      "moudjahid",
      "متحف المجاهد",
      "history",
      "independence",
    ],
    popularity: 6,
  },

  // Restaurants & Entertainment
  {
    id: "restaurant-el-dey",
    name: "Restaurant El Dey",
    arabicName: "مطعم الداي",
    category: "Restaurant",
    lat: 36.7538,
    lng: 3.0588,
    district: "Alger Centre",
    keywords: ["restaurant", "مطعم", "food", "dining", "el dey"],
    popularity: 6,
  },
  {
    id: "le-dauphin",
    name: "Le Dauphin",
    arabicName: "الدلفين",
    category: "Restaurant",
    lat: 36.7667,
    lng: 3.0333,
    district: "Alger Centre",
    keywords: ["restaurant", "le dauphin", "french", "dining", "seafood"],
    popularity: 7,
  },

  // Sports & Recreation
  {
    id: "stade-5-juillet",
    name: "Stade 5 Juillet",
    arabicName: "ملعب 5 جويلية",
    category: "Stadium",
    lat: 36.7403,
    lng: 3.0508,
    district: "Alger Centre",
    keywords: ["stadium", "5 juillet", "ملعب", "sports", "football", "soccer"],
    popularity: 8,
  },
  {
    id: "palais-des-sports",
    name: "Palais des Sports",
    arabicName: "قصر الرياضة",
    category: "Sports Center",
    lat: 36.7538,
    lng: 3.0588,
    district: "Alger Centre",
    keywords: ["sports", "palais", "قصر الرياضة", "indoor", "basketball"],
    popularity: 6,
  },

  // Religious Sites
  {
    id: "djamaa-el-djazair",
    name: "Djamaa El Djazair",
    arabicName: "جامع الجزائر",
    category: "Mosque",
    lat: 36.7394,
    lng: 3.1153,
    district: "Mohammadia",
    keywords: ["mosque", "djamaa", "جامع الجزائر", "great mosque", "islamic"],
    popularity: 9,
  },
  {
    id: "mosquee-emir-abdelkader",
    name: "Mosquée Emir Abdelkader",
    arabicName: "مسجد الأمير عبد القادر",
    category: "Mosque",
    lat: 36.7333,
    lng: 3.0833,
    district: "Kouba",
    keywords: ["mosque", "emir abdelkader", "مسجد", "islamic", "prayer"],
    popularity: 7,
  },
  {
    id: "basilique-notre-dame",
    name: "Basilique Notre-Dame d'Afrique",
    arabicName: "كنيسة سيدة إفريقيا",
    category: "Church",
    lat: 36.7908,
    lng: 3.0597,
    district: "Bologhine",
    keywords: ["church", "basilique", "notre dame", "christian", "catholic"],
    popularity: 7,
  },

  // Business Districts
  {
    id: "pins-maritimes",
    name: "Pins Maritimes",
    arabicName: "الصنوبر البحري",
    category: "Business District",
    lat: 36.7833,
    lng: 2.9167,
    district: "Sidi Fredj",
    keywords: ["business", "pins maritimes", "commercial", "office"],
    popularity: 6,
  },
  {
    id: "cyber-parc",
    name: "Cyber Parc Sidi Abdallah",
    arabicName: "الحديقة الإلكترونية سيدي عبد الله",
    category: "Technology Park",
    lat: 36.6833,
    lng: 2.8833,
    district: "Sidi Abdallah",
    keywords: ["cyber parc", "technology", "business", "IT", "innovation"],
    popularity: 6,
  },
];

export const searchAlgiersPlaces = (
  query: string,
  limit = 10
): AlgiersPlace[] => {
  if (!query.trim()) return [];

  const searchTerm = query.toLowerCase().trim();

  return algiersPlaces
    .filter((place) => {
      // Search in name, arabic name, and keywords
      const searchableText = [
        place.name.toLowerCase(),
        place.arabicName?.toLowerCase() || "",
        place.category.toLowerCase(),
        place.district.toLowerCase(),
        ...place.keywords.map((k) => k.toLowerCase()),
      ].join(" ");

      return searchableText.includes(searchTerm);
    })
    .sort((a, b) => {
      // Sort by relevance and popularity
      const aStartsWithQuery = a.name.toLowerCase().startsWith(searchTerm);
      const bStartsWithQuery = b.name.toLowerCase().startsWith(searchTerm);

      if (aStartsWithQuery && !bStartsWithQuery) return -1;
      if (!aStartsWithQuery && bStartsWithQuery) return 1;

      // If both start with query or neither does, sort by popularity
      return b.popularity - a.popularity;
    })
    .slice(0, limit);
};

// Helper function to get places by category
export const getPlacesByCategory = (category: string): AlgiersPlace[] => {
  return algiersPlaces.filter(
    (place) => place.category.toLowerCase() === category.toLowerCase()
  );
};

// Helper function to get popular places
export const getPopularPlaces = (limit = 10): AlgiersPlace[] => {
  return algiersPlaces
    .sort((a, b) => b.popularity - a.popularity)
    .slice(0, limit);
};

// Helper function to get places by district
export const getPlacesByDistrict = (district: string): AlgiersPlace[] => {
  return algiersPlaces.filter((place) =>
    place.district.toLowerCase().includes(district.toLowerCase())
  );
};
