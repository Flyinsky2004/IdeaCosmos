import {saveAs} from "file-saver";
import {AlignmentType, Document, HeadingLevel, ImageRun, Packer, Paragraph, TextRun} from "docx";
import html2pdf from "html2pdf.js";
import {get} from "@/util/request.js";
import {message} from "ant-design-vue";
import {imagePrefix} from "@/util/VARRIBLES.js";

export {exportScript, downloadWebp}
const downloadWebp = async (filename) => {
    let result;
    await new Promise((resolve, reject) => {
        get('/api/user/getWebpImageBase64', {
            filename: filename,
        }, (messager, data) => {
            result = data;
            resolve()
        }, (messager, data) => {
            message.warn(messager)
        }, (messager, data) => {
            message.error(messager)
        })
    })
    return result;
};

// 将 ArrayBuffer 转换为 Base64 的工具函数
const arrayBufferToBase64 = (buffer) => {
    let binary = '';
    const bytes = new Uint8Array(buffer);
    for (let i = 0; i < bytes.byteLength; i++) {
        binary += String.fromCharCode(bytes[i]);
    }
    return btoa(binary);
};

function base64ToUint8Array(base64) {
    const binaryString = atob(base64);
    const length = binaryString.length;
    const uint8Array = new Uint8Array(length);
    for (let i = 0; i < length; i++) {
        uint8Array[i] = binaryString.charCodeAt(i);
    }
    return uint8Array;
}

// 修改后的PDF导出函数
const exportToPDF = async (project, title, chapterTitle, chapterContent) => {
    try {
        const coverBase64 = await downloadWebp(project.cover_image)

        let htmlContent = `
            <div style="page-break-after: always;">
                <div style="text-align: center; margin-top: 50px;">
                    <img src="${coverBase64}" style="max-width: 100%; height: auto;">
                </div>

                <!-- 项目信息 -->
                <div style="margin: 30px 0;">
                    <h1 style="font-size: 24px; border-left: 4px solid #1890ff; padding-left: 12px;">项目信息</h1>
                    <ul style="list-style: none; padding: 0;">
                        <li><strong>项目名称：</strong>${project.project_name}</li>
                        <li><strong>作品类型：</strong>${project.types}</li>
                        <li><strong>作品风格：</strong>${project.style.join('、')}</li>
                        <li><strong>目标人群：</strong>${project.market_people.join('、')}</li>
                        <li><strong>剧情简述：</strong>${project.social_story}${project.start}</li>
                    </ul>
                </div>

                <!-- 团队信息 -->
                <div style="margin: 30px 0;">
                    <h2 style="font-size: 20px; border-left: 3px solid #666; padding-left: 10px;">创作团队</h2>
                    <div style="background: #f5f5f5; padding: 15px; border-radius: 4px;">
                        <p><strong>${project.team.username}</strong></p>
                        <p>${project.team.teamDescription}</p>
                    </div>
                </div>
            </div>

            <!-- 内容主体 -->
            <div style="page-break-before: always;">
                <h1 style="font-size: 28px; text-align: center;">${title}</h1>
                <h2 style="font-size: 22px; color: #666; text-align: center;">${chapterTitle}</h2>
                <div style="line-height: 1.8; font-size: 16px; margin: 20px 0;">
                    ${chapterContent.replace(/\n/g, "<br>")}
                </div>
            </div>
        `;

        const element = document.createElement("div");
        element.innerHTML = htmlContent;

        // 生成PDF
        await html2pdf()
            .from(element)
            .set({
                margin: [20, 15],
                filename: `${project.project_name} - ${chapterTitle}.pdf`,
                image: {type: 'jpeg', quality: 0.95},
                html2canvas: {
                    scale: 2,
                    logging: true,
                    useCORS: true,
                    allowTaint: false,
                    backgroundColor: '#FFFFFF'
                },
                jsPDF: {
                    unit: 'mm',
                    format: 'a4',
                    orientation: 'portrait',
                    compress: true
                },
                pagebreak: {
                    mode: 'avoid',
                    avoidBreakAfter: '.chapter-content',
                    avoidBreakBefore: '.chapter-content'
                }
            })
            .save();

    } catch (error) {
        console.error('导出PDF失败:', error);
        alert('导出失败，请检查网络连接或联系支持人员');
    }
};

const exportToWord = async (project, title, chapterTitle, chapterContent) => {
    const coverBase64 = await downloadWebp(project.cover_image)
    // 移除前缀（如果有）
    const base64Data = coverBase64.replace(/^data:image\/\w+;base64,/, "");
    const imageData = base64ToUint8Array(base64Data);

    const doc = new Document({
        sections: [{
            properties: {},
            children: [
                // 封面页
                new Paragraph({
                    children: [
                        new ImageRun({
                            data: imageData, // 需要替换为实际图片数据
                            transformation: {
                                width: 400,
                                height: 600
                            }
                        })
                    ],
                    alignment: AlignmentType.CENTER
                }),
                new Paragraph({
                    text: "项目信息",
                    heading: HeadingLevel.HEADING_1,
                    border: {left: {color: "1890ff", size: 8}}, // 左侧蓝色边框
                    spacing: {before: 800, after: 400}
                }),
                new Paragraph({
                    children: [
                        new TextRun({
                            text: "项目名称：",
                            bold: true
                        }),
                        new TextRun(project.project_name)
                    ]
                }),
                new Paragraph({
                    children: [
                        new TextRun({
                            text: "作品类型：",
                            bold: true
                        }),
                        new TextRun(project.types)
                    ]
                }),
                new Paragraph({
                    children: [
                        new TextRun({
                            text: "作品风格：",
                            bold: true
                        }),
                        new TextRun(project.style.join('、'))
                    ]
                }),
                new Paragraph({
                    children: [
                        new TextRun({
                            text: "目标人群：",
                            bold: true
                        }),
                        new TextRun(project.market_people.join('、'))
                    ]
                }),
                new Paragraph({
                    children: [
                        new TextRun({
                            text: "剧情简述：",
                            bold: true
                        }),
                        new TextRun(project.social_story + project.start)
                    ]
                }),

                // 团队信息
                new Paragraph({
                    text: "创作团队",
                    heading: HeadingLevel.HEADING_2,
                    border: {left: {color: "666666", size: 6}},
                    spacing: {before: 800, after: 400}
                }),
                new Paragraph({
                    children: [
                        new TextRun({
                            text: project.team.username,
                            bold: true
                        })
                    ],
                    shading: {fill: "F5F5F5"}, // 灰色背景
                    indent: {left: 400},
                    spacing: {line: 360}
                }),
                new Paragraph({
                    text: project.team.teamDescription,
                    shading: {fill: "F5F5F5"},
                    indent: {left: 400},
                    spacing: {line: 360}
                }),

                // 内容主体
                new Paragraph({
                    text: title,
                    heading: HeadingLevel.TITLE,
                    alignment: AlignmentType.CENTER,
                    pageBreakBefore: true
                }),
                new Paragraph({
                    text: chapterTitle,
                    heading: HeadingLevel.HEADING_1,
                    alignment: AlignmentType.CENTER,
                    spacing: {before: 600}
                }),
                new Paragraph({
                    text: chapterContent,
                    indent: {firstLine: 720}, // 首行缩进
                    spacing: {line: 380} // 行间距
                })
            ]
        }]
    });

    const blob = await Packer.toBlob(doc);
    saveAs(blob, `${project.project_name} - ${chapterTitle}.docx`);
};

const exportToMarkdown = (project, title, chapterTitle, chapterContent) => {
    let markdownContent = `![封面图片](` + imagePrefix + project.cover_image + `)\n\n`;

    // 项目信息
    markdownContent += "## 项目信息\n";
    markdownContent += `**项目名称：** ${project.project_name}\n\n`;
    markdownContent += `**作品类型：** ${project.types}\n\n`;
    markdownContent += `**作品风格：** ${project.style.join('、')}\n\n`;
    markdownContent += `**目标人群：** ${project.market_people.join('、')}\n\n`;
    markdownContent += `**剧情简述：** ${project.social_story}${project.start}\n\n`;

    // 团队信息
    markdownContent += "## 创作团队\n";
    markdownContent += `**${project.team.username}**\n\n`;
    markdownContent += `${project.team.teamDescription}\n\n`;

    // 内容主体
    markdownContent += `# ${title}\n\n`;
    markdownContent += `## ${chapterTitle}\n\n`;
    markdownContent += `${chapterContent.replace(/\n/g, "\n\n")}\n`;

    const blob = new Blob([markdownContent], {type: "text/markdown;charset=utf-8"});
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${project.project_name} - ${chapterTitle}.md`;
    a.style.display = "none";
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
};


const exportScript = (project, title, chapterTitle, chapterContent, format) => {
    const afterContent = chapterContent.replace(/\n/g, '\n  ');
    if (format === "pdf") {
        exportToPDF(project, title, chapterTitle, afterContent);
    } else if (format === "word") {
        exportToWord(project, title, chapterTitle, afterContent);
    } else if (format === "markdown") {
        exportToMarkdown(project, title, chapterTitle, afterContent);
    } else {
        console.error("Unsupported format");
    }
};