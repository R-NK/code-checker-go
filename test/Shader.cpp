#include "stdafx.h"
#include "Shader.h"
#include <iostream>
#include <fstream>
#include <sstream>
#include <glm/glm.hpp>
#include <glm/gtc/type_ptr.hpp>

Shader::Shader(const GLchar * vertexPath, const GLchar * fragmentPath)
{
	// シェーダのソースを読み込む
	std::string vsrcStr;
	readShaderSource(vertexPath, vsrcStr);
	auto vsrc = vsrcStr.data();

	std::string fsrcStr;
	readShaderSource(fragmentPath, fsrcStr);
	auto fsrc = fsrcStr.data();

	int_fast16_t success;
	GLchar infoLog[512];

	// 空のプログラムオブジェクトを作成する
	ID = glCreateProgram();

	if (vsrc != nullptr)
	{
		// バーテックスシェーダのシェーダオブジェクトを作成する
		const auto vobj(glCreateShader(GL_VERTEX_SHADER));
		glShaderSource(vobj, 1, &vsrc, nullptr);
		glCompileShader(vobj);

		// コンパイルチェック
		glGetShaderiv(vobj, GL_COMPILE_STATUS, &success);
		if (!success)
		{
			glGetShaderInfoLog(vobj, 512, nullptr, infoLog);
			std::cout << "ERROR::Shader::Vertex::Compilation_Failed\n" << infoLog << std::endl;
		}
		else
		{
			std::cout << "SUCCESS::Shader::Vertex::Compilation_Success\n" << std::endl;
		}

		// バーテックスシェーダのシェーダオブジェクトをプログラムオブジェクトに組み込む
		glAttachShader(ID, vobj);
		glDeleteShader(vobj);
	}

	if (fsrc != nullptr)
	{
		// フラグメントシェーダのシェーダオブジェクトを作成する
		const auto fobj(glCreateShader(GL_FRAGMENT_SHADER));
		glShaderSource(fobj, 1, &fsrc, nullptr);
		glCompileShader(fobj);

		// コンパイルチェック
		glGetShaderiv(fobj, GL_COMPILE_STATUS, &success);
		if (!success)
		{
			glGetShaderInfoLog(fobj, 512, nullptr, infoLog);
			std::cout << "ERROR::Shader::Fragment::Compilation_Failed\n" << infoLog << std::endl;
		}
		else
		{
			std::cout << "SUCCESS::Shader::Fragment::Compilation_Success\n" << std::endl;
		}

		// フラグメントシェーダのシェーダオブジェクトをプログラムオブジェクトに組み込む
		glAttachShader(ID, fobj);
		glDeleteShader(fobj);
	}

	// プログラムオブジェクトをリンクする
	glLinkProgram(ID);
}

Shader::~Shader() = default;

bool Shader::readShaderSource(const GLchar * path, std::string & buffer)
{
	// pathのNULLチェック
	if (path == nullptr)
	{
		std::cout << "ERROR::1st argument is nullptr.\n" << std::endl;
		return false;
	}
	std::ifstream fin(path, std::ios::in);
	if (!fin)
	{
		std::cout << "ERROR::Can't read the file::" << path << std::endl;
		return false;
	}

	std::stringstream stringstream;
	stringstream << fin.rdbuf();
	fin.close();

	buffer = stringstream.str();

	return true;
}

void Shader::use() const
{
	glUseProgram(ID);
}

void Shader::setAttributeValue(const std::string & name, bool value) const
{
	glUniform1i(glGetUniformLocation(ID, name.c_str()), static_cast<GLint>(value));
}

void Shader::setAttributeValue(const std::string & name, int value) const
{
	glUniform1i(glGetUniformLocation(ID, name.c_str()), value);
}

void Shader::setAttributeValue(const std::string & name, float value) const
{
	glUniform1f(glGetUniformLocation(ID, name.c_str()), value);
}

void Shader::setAttributeValue(const std::string & name, const glm::mat4 & value) const
{
    glUniformMatrix4fv(glGetUniformLocation(ID, name.c_str()), 1, GL_FALSE, glm::value_ptr(value));
}
