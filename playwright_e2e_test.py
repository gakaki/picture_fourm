#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Nano Banana AI 图片生成系统 - Playwright 端到端测试
测试前端界面功能和后端API集成
"""

import asyncio
import aiohttp
import json
import time
import os
from playwright.async_api import async_playwright
import sys

class NanoBananaE2ETest:
    def __init__(self):
        self.frontend_url = "http://localhost:3000"
        self.backend_url = "http://localhost:8080"
        self.test_results = []
        self.failed_tests = []
        
    async def log_test(self, test_name, status, details=""):
        """记录测试结果"""
        result = {
            "test": test_name,
            "status": status,
            "details": details,
            "timestamp": time.strftime("%Y-%m-%d %H:%M:%S")
        }
        self.test_results.append(result)
        
        status_icon = "✅" if status == "PASS" else "❌"
        print(f"{status_icon} {test_name}: {status}")
        if details:
            print(f"   详情: {details}")
        
        if status == "FAIL":
            self.failed_tests.append(test_name)

    async def test_backend_health(self):
        """测试后端健康检查"""
        try:
            async with aiohttp.ClientSession() as session:
                async with session.get(f"{self.backend_url}/api/v1/health") as resp:
                    if resp.status == 200:
                        data = await resp.json()
                        await self.log_test(
                            "后端健康检查", 
                            "PASS", 
                            f"状态码: {resp.status}, 响应: {data}"
                        )
                        return True
                    else:
                        await self.log_test(
                            "后端健康检查", 
                            "FAIL", 
                            f"状态码: {resp.status}"
                        )
                        return False
        except Exception as e:
            await self.log_test("后端健康检查", "FAIL", f"异常: {str(e)}")
            return False

    async def test_backend_text2img_api(self):
        """测试后端文本生成图片API"""
        try:
            test_payload = {
                "prompt": "一只可爱的小猫咪，蓝色的眼睛",
                "count": 1,
                "params": {"size": "512x512", "quality": "standard"}
            }
            
            async with aiohttp.ClientSession() as session:
                async with session.post(
                    f"{self.backend_url}/api/v1/generate/text2img",
                    json=test_payload,
                    headers={"Content-Type": "application/json"}
                ) as resp:
                    if resp.status == 200:
                        data = await resp.json()
                        if data.get("success") and data.get("data"):
                            image_url = data["data"][0].get("image_url")
                            generation_time = data["data"][0].get("generation_time", 0)
                            await self.log_test(
                                "后端文本生成图片API", 
                                "PASS", 
                                f"图片URL: {image_url}, 生成时间: {generation_time:.2f}s"
                            )
                            return True
                        else:
                            await self.log_test(
                                "后端文本生成图片API", 
                                "FAIL", 
                                f"响应格式错误: {data}"
                            )
                            return False
                    else:
                        text = await resp.text()
                        await self.log_test(
                            "后端文本生成图片API", 
                            "FAIL", 
                            f"状态码: {resp.status}, 响应: {text}"
                        )
                        return False
        except Exception as e:
            await self.log_test("后端文本生成图片API", "FAIL", f"异常: {str(e)}")
            return False

    async def test_generated_image_access(self):
        """测试生成图片的访问"""
        try:
            # 先生成一张图片
            test_payload = {
                "prompt": "测试图片访问",
                "count": 1,
                "params": {"size": "512x512"}
            }
            
            async with aiohttp.ClientSession() as session:
                # 生成图片
                async with session.post(
                    f"{self.backend_url}/api/v1/generate/text2img",
                    json=test_payload
                ) as resp:
                    if resp.status != 200:
                        await self.log_test("图片访问测试", "FAIL", "图片生成失败")
                        return False
                    
                    data = await resp.json()
                    image_url = data["data"][0]["image_url"]
                    full_image_url = f"{self.backend_url}{image_url}"
                    
                # 访问图片
                async with session.get(full_image_url) as img_resp:
                    if img_resp.status == 200:
                        content_type = img_resp.headers.get('Content-Type', '')
                        content_length = img_resp.headers.get('Content-Length', '0')
                        await self.log_test(
                            "图片访问测试", 
                            "PASS", 
                            f"Content-Type: {content_type}, Size: {content_length} bytes"
                        )
                        return True
                    else:
                        await self.log_test(
                            "图片访问测试", 
                            "FAIL", 
                            f"图片访问失败，状态码: {img_resp.status}"
                        )
                        return False
        except Exception as e:
            await self.log_test("图片访问测试", "FAIL", f"异常: {str(e)}")
            return False

    async def test_frontend_loading(self, page):
        """测试前端页面加载"""
        try:
            await page.goto(self.frontend_url)
            await page.wait_for_load_state('networkidle')
            
            # 检查页面标题
            title = await page.title()
            if "Nano Banana" in title or title:
                await self.log_test("前端页面加载", "PASS", f"页面标题: {title}")
                return True
            else:
                await self.log_test("前端页面加载", "FAIL", f"页面标题异常: {title}")
                return False
        except Exception as e:
            await self.log_test("前端页面加载", "FAIL", f"异常: {str(e)}")
            return False

    async def test_frontend_navigation(self, page):
        """测试前端页面导航"""
        try:
            await page.goto(self.frontend_url)
            await page.wait_for_load_state('networkidle')
            
            # 查找导航链接或按钮
            navigation_elements = await page.query_selector_all('nav a, [role="navigation"] a, .nav-link')
            
            if len(navigation_elements) > 0:
                await self.log_test(
                    "前端页面导航", 
                    "PASS", 
                    f"找到 {len(navigation_elements)} 个导航元素"
                )
                return True
            else:
                # 检查是否是单页应用
                main_content = await page.query_selector('main, .main, #root, .app')
                if main_content:
                    await self.log_test("前端页面导航", "PASS", "单页应用结构正常")
                    return True
                else:
                    await self.log_test("前端页面导航", "FAIL", "未找到导航元素或主要内容区域")
                    return False
        except Exception as e:
            await self.log_test("前端页面导航", "FAIL", f"异常: {str(e)}")
            return False

    async def test_image_generation_ui(self, page):
        """测试前端图片生成界面"""
        try:
            await page.goto(self.frontend_url)
            await page.wait_for_load_state('networkidle')
            
            # 查找文本输入框
            text_inputs = await page.query_selector_all('input[type="text"], textarea, [placeholder*="提示"], [placeholder*="prompt"]')
            
            # 查找生成按钮
            generate_buttons = await page.query_selector_all('button:has-text("生成"), button:has-text("Generate"), [type="submit"]')
            
            if len(text_inputs) > 0 and len(generate_buttons) > 0:
                await self.log_test(
                    "图片生成界面", 
                    "PASS", 
                    f"找到 {len(text_inputs)} 个输入框和 {len(generate_buttons)} 个按钮"
                )
                return True
            else:
                await self.log_test(
                    "图片生成界面", 
                    "FAIL", 
                    f"输入框: {len(text_inputs)}, 按钮: {len(generate_buttons)}"
                )
                return False
        except Exception as e:
            await self.log_test("图片生成界面", "FAIL", f"异常: {str(e)}")
            return False

    async def test_frontend_image_generation(self, page):
        """测试前端完整图片生成流程"""
        try:
            await page.goto(self.frontend_url)
            await page.wait_for_load_state('networkidle')
            
            # 尝试找到提示词输入框
            prompt_input = None
            for selector in [
                'input[placeholder*="提示"]',
                'input[placeholder*="prompt"]', 
                'textarea[placeholder*="提示"]',
                'textarea[placeholder*="prompt"]',
                'input[type="text"]',
                'textarea'
            ]:
                element = await page.query_selector(selector)
                if element:
                    prompt_input = element
                    break
            
            if not prompt_input:
                await self.log_test("前端图片生成流程", "FAIL", "未找到提示词输入框")
                return False
            
            # 输入测试提示词
            test_prompt = "一只可爱的小猫咪，蓝色的眼睛"
            await prompt_input.fill(test_prompt)
            
            # 查找并点击生成按钮
            generate_button = None
            for selector in [
                'button:has-text("生成")',
                'button:has-text("Generate")',
                'button[type="submit"]',
                '.generate-btn',
                '#generate-btn'
            ]:
                element = await page.query_selector(selector)
                if element:
                    generate_button = element
                    break
            
            if not generate_button:
                await self.log_test("前端图片生成流程", "FAIL", "未找到生成按钮")
                return False
            
            # 点击生成按钮
            await generate_button.click()
            
            # 等待请求完成（最长等待30秒）
            try:
                # 监听网络响应
                response = None
                def handle_response(resp):
                    nonlocal response
                    if "/api/v1/generate/text2img" in resp.url:
                        response = resp
                
                page.on("response", handle_response)
                
                # 等待最长30秒获取响应
                for _ in range(30):
                    if response:
                        break
                    await page.wait_for_timeout(1000)  # 等待1秒
                
                if response and response.status < 400:
                    await self.log_test("前端图片生成流程", "PASS", f"成功发送图片生成请求，状态码: {response.status}")
                    return True
                else:
                    status = response.status if response else "无响应"
                    await self.log_test(
                        "前端图片生成流程", 
                        "FAIL", 
                        f"API响应异常，状态码: {status}"
                    )
                    return False
            except Exception as timeout_e:
                await self.log_test(
                    "前端图片生成流程", 
                    "FAIL", 
                    f"等待API响应异常: {str(timeout_e)}"
                )
                return False
                
        except Exception as e:
            await self.log_test("前端图片生成流程", "FAIL", f"异常: {str(e)}")
            return False

    async def test_prompts_api(self):
        """测试提示词管理API"""
        try:
            async with aiohttp.ClientSession() as session:
                # 测试获取提示词列表
                async with session.get(f"{self.backend_url}/api/v1/prompts") as resp:
                    if resp.status == 200:
                        data = await resp.json()
                        await self.log_test(
                            "提示词管理API", 
                            "PASS", 
                            f"成功获取提示词列表，数量: {len(data.get('data', {}).get('prompts', []))}"
                        )
                        return True
                    else:
                        text = await resp.text()
                        await self.log_test(
                            "提示词管理API", 
                            "FAIL", 
                            f"状态码: {resp.status}, 响应: {text}"
                        )
                        return False
        except Exception as e:
            await self.log_test("提示词管理API", "FAIL", f"异常: {str(e)}")
            return False

    async def run_all_tests(self):
        """运行所有测试"""
        print("🚀 开始 Nano Banana AI 图片生成系统端到端测试")
        print("=" * 60)
        
        # 后端API测试
        print("\n📡 后端API测试")
        print("-" * 30)
        await self.test_backend_health()
        await self.test_backend_text2img_api()
        await self.test_generated_image_access()
        await self.test_prompts_api()
        
        # 前端UI测试
        print("\n🎨 前端UI测试")
        print("-" * 30)
        
        async with async_playwright() as p:
            browser = await p.chromium.launch(headless=True)
            context = await browser.new_context()
            page = await context.new_page()
            
            await self.test_frontend_loading(page)
            await self.test_frontend_navigation(page)
            await self.test_image_generation_ui(page)
            await self.test_frontend_image_generation(page)
            
            await browser.close()
        
        # 生成测试报告
        await self.generate_report()

    async def generate_report(self):
        """生成测试报告"""
        print("\n" + "=" * 60)
        print("📊 测试报告")
        print("=" * 60)
        
        total_tests = len(self.test_results)
        passed_tests = len([r for r in self.test_results if r["status"] == "PASS"])
        failed_tests = len(self.failed_tests)
        
        print(f"总测试数: {total_tests}")
        print(f"通过: {passed_tests}")
        print(f"失败: {failed_tests}")
        print(f"通过率: {(passed_tests/total_tests)*100:.1f}%")
        
        if self.failed_tests:
            print(f"\n❌ 失败的测试:")
            for test in self.failed_tests:
                print(f"   - {test}")
        
        # 保存详细报告到文件
        report_file = "test_report.json"
        report_data = {
            "summary": {
                "total": total_tests,
                "passed": passed_tests,
                "failed": failed_tests,
                "pass_rate": f"{(passed_tests/total_tests)*100:.1f}%"
            },
            "failed_tests": self.failed_tests,
            "detailed_results": self.test_results
        }
        
        with open(report_file, "w", encoding="utf-8") as f:
            json.dump(report_data, f, ensure_ascii=False, indent=2)
        
        print(f"\n📄 详细报告已保存到: {report_file}")
        
        if failed_tests == 0:
            print("\n🎉 所有测试通过！系统运行正常。")
        else:
            print(f"\n⚠️  发现 {failed_tests} 个问题需要修复。")

async def main():
    """主函数"""
    tester = NanoBananaE2ETest()
    await tester.run_all_tests()

if __name__ == "__main__":
    asyncio.run(main())