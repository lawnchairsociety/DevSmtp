using DevSmtp.Core.Utilities;
using Microsoft.VisualStudio.TestTools.UnitTesting;

namespace DevSmtp.Core.Tests.Utilities
{
    [TestClass]
    public class StringExtensionsTests
    {
        [TestMethod]
        public void IsEmpty_WhenStringIsValid_ItShouldReturnFalse()
        {
            // Arrange
            var candidate = "sometext";

            // Act
            var results = candidate.IsEmpty();

            // Assert
            Assert.IsFalse(results);
        }

        [TestMethod]
        public void IsEmpty_WhenStringIsNull_ItShouldReturnTrue()
        {
            // Arrange
            string? candidate = null;

            // Act
            var results = candidate.IsEmpty();

            // Assert
            Assert.IsTrue(results);
        }

        [TestMethod]
        public void IsEmpty_WhenStringIsWhitespace_ItShouldReturnTrue()
        {
            // Arrange
            var candidate = " ";

            // Act
            var results = candidate.IsEmpty();

            // Assert
            Assert.IsTrue(results);
        }
    }
}
