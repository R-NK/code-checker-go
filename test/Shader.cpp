#include "stdafx.h"
#include "Shader.h"
#include <iostream>
#include <fstream>
#include <sstream>
#include <glm/glm.hpp>
#include <glm/gtc/type_ptr.hpp>

Shader::Shader(const GLchar * vertexPath, const GLchar * fragmentPath)
{
	// �V�F�[�_�̃\�[�X��ǂݍ���
	std::string vsrcStr;
	readShaderSource(vertexPath, vsrcStr);
	auto vsrc = vsrcStr.data();

	std::string fsrcStr;
	readShaderSource(fragmentPath, fsrcStr);
	auto fsrc = fsrcStr.data();

	int_fast16_t success;
	GLchar infoLog[512];

	// ��̃v���O�����I�u�W�F�N�g���쐬����
	ID = glCreateProgram();

	if (vsrc != nullptr)
	{
		// �o�[�e�b�N�X�V�F�[�_�̃V�F�[�_�I�u�W�F�N�g���쐬����
		const auto vobj(glCreateShader(GL_VERTEX_SHADER));
		glShaderSource(vobj, 1, &vsrc, nullptr);
		glCompileShader(vobj);

		// �R���p�C���`�F�b�N
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

		// �o�[�e�b�N�X�V�F�[�_�̃V�F�[�_�I�u�W�F�N�g���v���O�����I�u�W�F�N�g�ɑg�ݍ���
		glAttachShader(ID, vobj);
		glDeleteShader(vobj);
	}

	if (fsrc != nullptr)
	{
		// �t���O�����g�V�F�[�_�̃V�F�[�_�I�u�W�F�N�g���쐬����
		const auto fobj(glCreateShader(GL_FRAGMENT_SHADER));
		glShaderSource(fobj, 1, &fsrc, nullptr);
		glCompileShader(fobj);

		// �R���p�C���`�F�b�N
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

		// �t���O�����g�V�F�[�_�̃V�F�[�_�I�u�W�F�N�g���v���O�����I�u�W�F�N�g�ɑg�ݍ���
		glAttachShader(ID, fobj);
		glDeleteShader(fobj);
	}

	// �v���O�����I�u�W�F�N�g�������N����
	glLinkProgram(ID);
}

Shader::~Shader() = default;

bool Shader::readShaderSource(const GLchar * path, std::string & buffer)
{
	// path��NULL�`�F�b�N
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
