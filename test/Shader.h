#pragma once
#include <GL/glew.h>
#include <string>
#include <glm/glm.hpp>

class Shader
{
public:
	// program ID
	unsigned int ID;

	Shader(const GLchar* vertexPath, const GLchar* fragmentPath);
	// シェ&#65533;[ダを有効化
	void use() const;
	// utility uniform functions
	void setAttributeValue(const std::string &name, bool value) const;
	void setAttributeValue(const std::string &name, int value) const;
	void setAttributeValue(const std::string &name, float value) const;
    void setAttributeValue(const std::string &name, const glm::mat4& value) const;
	~Shader();
private:
	static bool readShaderSource(const GLchar* path, std::string& buffer);
};
