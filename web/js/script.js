function print_doc(){
    window.print();
}
     
function print_tabdoc() {
    var docDefinition = {
content: [
{ text: 'Документы', style: 'header' },
'Официальная документация по работе с сервисом OrionWeb',
{ text: 'Справочник пользователей', style: 'subheader' },
'Представление таблицы пользователей',
{
style: 'tableExample',
table: {
    body: [
        ['name', 'password'],
        ['Alex', 'c4ca4238a0b923820dcc509a6f75849b']
    ]
}
},
{ text: 'Справочник организации', style: 'subheader' },
'Представление таблицы организации',
{
style: 'tableExample',
table: {
    body: [
        ['id', 'name', 'inn'],
        ['1', 'ООО БД Безопасность', '702148941824']
    ]
}
},
{ text: 'Справочник сотрудников', style: 'subheader' },
'Представление таблицы сотрудников',
{
    table: {
        body: [
            ['surname', 'name', 'patronymic', 'subdivision', 'organisation'],
            ['Сафронов', 'Алексей', '', '', '']
        ]
    }
},
{ text: 'Справочник подразделений', style: 'subheader' },
{
style: 'tableExample',
table: {
    body: [
        ['id', 'name', 'description'],
        ['1', 'Подразделение ИБ', '']
    ]
}
}
],
styles: {
header: {
fontSize: 18,
bold: true,
margin: [0, 0, 0, 10]
},
subheader: {
fontSize: 16,
bold: true,
margin: [0, 10, 0, 5]
},
tableExample: {
margin: [0, 5, 0, 15]
},
tableOpacityExample: {
margin: [0, 5, 0, 15],
fillColor: 'blue',
fillOpacity: 0.3
},
tableHeader: {
bold: true,
fontSize: 13,
color: 'black'
}
},
defaultStyle: {
// alignment: 'justify'
},
patterns: {
stripe45d: {
boundingBox: [1, 1, 4, 4],
xStep: 3,
yStep: 3,
pattern: '1 w 0 1 m 4 5 l s 2 0 m 5 3 l s'
}
}
};
    pdfMake.createPdf(docDefinition).open();
}