using DevSmtp.Core.Stores;

namespace DevSmtp.Core.Commands
{
    public sealed class SamlHandler : ICommandHandler<Saml, SamlResult>
    {
        private readonly IDataStore _dataStore;

        public SamlHandler(IDataStore dataStore)
        {
            this._dataStore = dataStore ?? throw new ArgumentNullException(nameof(dataStore));
        }

        public Task<SamlResult> ExecuteAsync(Saml command, CancellationToken cancellationToken = default)
        {
            throw new NotImplementedException();
        }
    }
}
