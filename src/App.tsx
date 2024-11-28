import SearchText from "./lib/components/SearchText/SearchText";
import SearchResult from "./lib/components/SearchResults/SearchResults";
import { Topic } from "./lib/stores/SemanticSearchStore/interface";

export default function App() {
  const results: Topic[] = [
    {
      title: 'วิสัยทัศน์',
      context: `
"ราชมงคลพระนคร"
มหาวิทยาลัยแห่งเทคโนโลยีนวัตกรรม
และการบูรณาการ
      `
    },
    {
      title: 'พันธกิจ',
      context: `
1.ผลิตและพัฒนากำลังคนให้พร้อมเป็น "นวัตกรบูรณาการ" ที่มีความรอบรู้ มีความสามารถในการปรับตัว และรับมือกับความท้าทายได้อย่างรวดเร็ว
2.สร้างฐานข้อมูลคุณภาพ เพื่อสร้างสรรค์งานวิจัย และต่อยอดนวัตกรรมใหม่ ๆ ที่ตอบสนองต่อความต้องการของสังคมและชุมชน
3.บริการวิชาการต่ออุตสาหกรรมเป้าหมาย เพื่อยกระดับการพัฒนาอย่างยั่งยืน
4.ทำนุบำรุงศาสนา ศิลปวัฒนธรรม และอนุรักษ์สิ่งแวดล้อม
5.บริหารจัดการอย่างมีธรรมาภิบาล พร้อมสร้างวัฒนธรรมองค์กรต้นแบบ
      `
    }
  ]

  return (
    <div className="container mx-auto">
      <div className="my-4">
        <SearchText />
        <SearchResult results={results} />
      </div>
    </div>
  )
}