using DevSmtp.Core.Models;
using Microsoft.VisualStudio.TestTools.UnitTesting;

namespace DevSmtp.Core.Tests.Models.Primitives
{
    [TestClass]
    public class EmailTests
    {
        [TestMethod]
        public void Ctor_WhenValueIsValid_ItShouldCreateEmail()
        {
            // Arrange
            var value = "user@fake.example.com";

            // Act
            var results = new Email(value);

            // Assert
            Assert.IsNotNull(results);
            Assert.AreEqual(value, results.Value);
        }

        [TestMethod]
        public void From_WhenValueIsValid_ItShouldCreateEmail()
        {
            // Arrange
            var value = "user@fake.example.com";

            // Act
            var results = Email.From(value);

            // Assert
            Assert.IsNotNull(results);
            Assert.AreEqual(value, results.Value);
        }

        [TestMethod]
        public void Ctor_WhenValueIsNotValid_ItShouldThrowFormatException()
        {
            // Arrange
            var value = "not_an_email_address";

            // Act
            try
            {
                _ = new Email(value);
                Assert.Fail("FormatException expected");
            }
            catch (FormatException)
            {
                // expected...
            }
        }
    }
}
